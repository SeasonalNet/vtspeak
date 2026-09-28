#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <windows.h>

_Static_assert(sizeof(uintptr_t) == 4, "PE32 module-address layout");

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
    char *record = malloc((size_t)file_length + 1U);
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

int main(int argc, char **argv) {
    if (argc != 2 || (strcmp(argv[1], "file") != 0 && strcmp(argv[1], "memory") != 0)) {
        puts("HARNESS usage=probe-load-ext-james-license.exe <file|memory>");
        return 2;
    }
    HMODULE module = LoadLibraryA("Z:\\samples\\vt_pau.dll");
    if (module == NULL) {
        printf("HARNESS load_library=%lu\n", (unsigned long)GetLastError());
        return 3;
    }
    load_ext_fn load_ext;
    check_license_fn check_license;
    get_db_size_fn get_db_size;
    unload_ext_fn unload_ext;
    FARPROC load_symbol = GetProcAddress(module, "VT_LOADTTS_EXT_ENG");
    FARPROC check_symbol = GetProcAddress(module, "VT_CheckLicense_ENG");
    FARPROC size_symbol = GetProcAddress(module, "VT_GetDBSize_ENG");
    FARPROC unload_symbol = GetProcAddress(module, "VT_UNLOADTTS_EXT_ENG");
    _Static_assert(sizeof(load_ext) == sizeof(load_symbol), "PE32 function pointer size");
    _Static_assert(sizeof(check_license) == sizeof(check_symbol), "PE32 function pointer size");
    _Static_assert(sizeof(get_db_size) == sizeof(size_symbol), "PE32 function pointer size");
    _Static_assert(sizeof(unload_ext) == sizeof(unload_symbol), "PE32 function pointer size");
    memcpy(&load_ext, &load_symbol, sizeof(load_ext));
    memcpy(&check_license, &check_symbol, sizeof(check_license));
    memcpy(&get_db_size, &size_symbol, sizeof(get_db_size));
    memcpy(&unload_ext, &unload_symbol, sizeof(unload_ext));
    if (load_ext == NULL || check_license == NULL || get_db_size == NULL || unload_ext == NULL) {
        puts("HARNESS missing_export=1");
        FreeLibrary(module);
        return 4;
    }
    char *record = NULL;
    uint32_t record_length = 0;
    if (!read_license_record(&record, &record_length)) {
        puts("HARNESS read_record=0");
        FreeLibrary(module);
        return 5;
    }

    char database_base[] = "Z:\\work\\";
    char license_path[] = "Z:\\work\\data-common\\verify\\verification.txt";
    int use_file = strcmp(argv[1], "file") == 0;
    short load_result = load_ext(NULL, 4, database_base, 0, UINT32_MAX,
                                 use_file ? license_path : NULL,
                                 use_file ? NULL : record,
                                 use_file ? UINT32_MAX : record_length);
    int database_size = 0;
    uint32_t size_result = get_db_size(&database_size, 4);
    uintptr_t base = (uintptr_t)module;
    const char *constraint =
        *(const char **)(base + 0x7c6a8U + 4U * 6U * sizeof(uintptr_t));
    int check_result = check_license(use_file ? license_path : NULL,
                                     use_file ? NULL : record,
                                     use_file ? UINT32_MAX : record_length,
                                     (unsigned char *)constraint);
    uintptr_t slot_state = *(uintptr_t *)(base + 0xa0464U + 4U * sizeof(uintptr_t));
    unsigned int gate = *(unsigned char *)(base + 0xa7488U + 4U);
    int capacity = slot_state == 0 ? -1 : *(int *)(slot_state + 0x4d14U);
    printf("JAMES_LICENSE form=%s bytes=%lu load_ax=%d dbsize_result=%lu dbsize=%d check=%d gate=%u capacity=%d\n",
           argv[1], (unsigned long)record_length, (int)load_result,
           (unsigned long)size_result, database_size, check_result, gate, capacity);

    if (slot_state != 0) {
        unload_ext(4);
    }
    free(record);
    FreeLibrary(module);
    return 0;
}
