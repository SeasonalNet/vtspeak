#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <windows.h>

typedef short(__cdecl *load_ext_fn)(HWND, int, char *, uint32_t, uint32_t, char *, char *,
                                    uint32_t);
typedef int(__cdecl *check_license_fn)(char *, char *, uint32_t, unsigned char *);
typedef uint32_t(__cdecl *get_db_size_fn)(int *, int);
typedef void(__cdecl *unload_ext_fn)(int);

static int read_license_record(char **record_out, uint32_t *length_out) {
    FILE *file = fopen("Z:\\work\\data-common\\verify\\verification.txt", "rb");
    if (file == NULL || fseek(file, 0, SEEK_END) != 0) {
        if (file != NULL) {
            fclose(file);
        }
        return 0;
    }
    long file_length = ftell(file);
    if (file_length <= 0 || fseek(file, 0, SEEK_SET) != 0) {
        fclose(file);
        return 0;
    }
    char *record = malloc((size_t)file_length + 1);
    if (record == NULL) {
        fclose(file);
        return 0;
    }
    size_t bytes_read = fread(record, 1, (size_t)file_length, file);
    fclose(file);
    if (bytes_read != (size_t)file_length) {
        free(record);
        return 0;
    }
    record[bytes_read] = '\0';
    *record_out = record;
    *length_out = (uint32_t)bytes_read;
    return 1;
}

int main(void) {
    HMODULE dll = LoadLibraryA("Z:\\samples\\vt_pau.dll");
    if (dll == NULL) {
        printf("HARNESS load_library=%lu\n", (unsigned long)GetLastError());
        return 2;
    }
    FARPROC load_symbol = GetProcAddress(dll, "VT_LOADTTS_EXT_ENG");
    FARPROC check_symbol = GetProcAddress(dll, "VT_CheckLicense_ENG");
    FARPROC size_symbol = GetProcAddress(dll, "VT_GetDBSize_ENG");
    FARPROC unload_symbol = GetProcAddress(dll, "VT_UNLOADTTS_EXT_ENG");
    load_ext_fn load_ext;
    check_license_fn check_license;
    get_db_size_fn get_db_size;
    unload_ext_fn unload_ext;
    _Static_assert(sizeof(load_ext) == sizeof(load_symbol), "PE32 function pointer size");
    _Static_assert(sizeof(check_license) == sizeof(check_symbol), "PE32 function pointer size");
    _Static_assert(sizeof(get_db_size) == sizeof(size_symbol), "PE32 function pointer size");
    _Static_assert(sizeof(unload_ext) == sizeof(unload_symbol), "PE32 function pointer size");
    memcpy(&load_ext, &load_symbol, sizeof(load_ext));
    memcpy(&check_license, &check_symbol, sizeof(check_license));
    memcpy(&get_db_size, &size_symbol, sizeof(get_db_size));
    memcpy(&unload_ext, &unload_symbol, sizeof(unload_ext));
    if (load_ext == NULL || check_license == NULL || get_db_size == NULL || unload_ext == NULL) {
        printf("HARNESS missing_export load=%d check=%d size=%d unload=%d\n", load_ext != NULL,
               check_license != NULL, get_db_size != NULL, unload_ext != NULL);
        FreeLibrary(dll);
        return 3;
    }
    char *record = NULL;
    uint32_t record_length = 0;
    if (!read_license_record(&record, &record_length)) {
        puts("HARNESS read_record=0");
        FreeLibrary(dll);
        return 4;
    }
    char database_base[] = "Z:\\work\\";
    short load_result = load_ext(NULL, 1, database_base, 0, UINT32_MAX, NULL, record,
                                 record_length);
    int database_size = 0;
    uint32_t size_result = get_db_size(&database_size, 1);
    uintptr_t module_base = (uintptr_t)dll;
    const char *voice_constraint =
        *(const char **)(module_base + 0x7c6a8U + 6U * sizeof(void *));
    int check_result = check_license(NULL, record, record_length,
                                     (unsigned char *)voice_constraint);
    int slot_state = *((int *)(module_base + 0xa0464U) + 1);
    int dictionary_capacity = *(int *)((uintptr_t)slot_state + 0x4d14U);
    unsigned int license_gate = *(unsigned char *)(module_base + 0xa7489U);
    printf("MEMORY_LICENSE bytes=%lu load_ax=%d dbsize_result=%lu dbsize=%d check=%d gate=%u capacity=%d\n",
           (unsigned long)record_length, (int)load_result, (unsigned long)size_result,
           database_size, check_result, license_gate, dictionary_capacity);
    if (size_result == 1) {
        unload_ext(1);
    }
    free(record);
    FreeLibrary(dll);
    return 0;
}
