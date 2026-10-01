#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <windows.h>

typedef int *(__cdecl *alloc_syncinfo_fn)(void);
typedef void(__cdecl *free_syncinfo_fn)(int *);

__attribute__((noinline)) void syncinfo_probe_alloc_start(void) {}
__attribute__((noinline)) void syncinfo_probe_alloc_done(int valid, int rows,
                                                          int width, int row_array,
                                                          int first_nested,
                                                          int last_nested,
                                                          int distinct_nested) {
    (void)valid;
    (void)rows;
    (void)width;
    (void)row_array;
    (void)first_nested;
    (void)last_nested;
    (void)distinct_nested;
}

int main(void) {
    HMODULE module = LoadLibraryA("Z:\\samples\\vt_pau.dll");
    if (module == NULL) {
        printf("SYNC_ALLOC_HARNESS load_library=%lu\n", (unsigned long)GetLastError());
        return 1;
    }

    FARPROC alloc_symbol = GetProcAddress(module, "VT_AllocSyncInfo_New_ENG");
    FARPROC free_symbol = GetProcAddress(module, "VT_FreeSyncInfo_New_ENG");
    alloc_syncinfo_fn alloc_syncinfo;
    free_syncinfo_fn free_syncinfo;
    _Static_assert(sizeof(alloc_syncinfo) == sizeof(alloc_symbol), "PE32 function pointer");
    _Static_assert(sizeof(free_syncinfo) == sizeof(free_symbol), "PE32 function pointer");
    memcpy(&alloc_syncinfo, &alloc_symbol, sizeof(alloc_syncinfo));
    memcpy(&free_syncinfo, &free_symbol, sizeof(free_syncinfo));
    if (alloc_syncinfo == NULL || free_syncinfo == NULL) {
        puts("SYNC_ALLOC_HARNESS missing_export=1");
        FreeLibrary(module);
        return 2;
    }

    syncinfo_probe_alloc_start();
    int *object = alloc_syncinfo();
    if (object == NULL) {
        puts("SYNC_ALLOC_RESULT object_null=1");
        FreeLibrary(module);
        return 3;
    }

    unsigned char *rows = (unsigned char *)(uintptr_t)(uint32_t)object[0];
    unsigned char *first_nested =
        (unsigned char *)(uintptr_t)*(uint32_t *)(rows + 4);
    unsigned char *last_nested =
        (unsigned char *)(uintptr_t)*(uint32_t *)(rows + 599U * 0x24U + 4U);
    int valid = object[1] == 600 && object[2] == 65 && rows != NULL &&
                first_nested != NULL && last_nested != NULL &&
                first_nested != last_nested;
    printf("SYNC_ALLOC_RESULT valid=%d rows=%d width=%d row_array=%d first_nested=%d last_nested=%d distinct_nested=%d\n",
           valid, object[1], object[2], rows != NULL, first_nested != NULL,
           last_nested != NULL, first_nested != last_nested);
    syncinfo_probe_alloc_done(valid, object[1], object[2], rows != NULL,
                              first_nested != NULL, last_nested != NULL,
                              first_nested != last_nested);

    free_syncinfo(object);
    FreeLibrary(module);
    return valid ? 0 : 4;
}
