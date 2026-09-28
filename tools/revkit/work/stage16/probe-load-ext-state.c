#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <windows.h>

_Static_assert(sizeof(uintptr_t) == 4, "PE32 module-address layout");

typedef short(__cdecl *load_ext_fn)(HWND, int, char *, uint32_t, uint32_t, char *, char *,
                                    uint32_t);
typedef uint32_t(__cdecl *get_db_size_fn)(int *, int);
typedef void(__cdecl *unload_ext_fn)(int);

static load_ext_fn resolve_load(HMODULE module) {
    FARPROC symbol = GetProcAddress(module, "VT_LOADTTS_EXT_ENG");
    load_ext_fn function;
    _Static_assert(sizeof(function) == sizeof(symbol), "PE32 function pointer size");
    memcpy(&function, &symbol, sizeof(function));
    return function;
}

static get_db_size_fn resolve_size(HMODULE module) {
    FARPROC symbol = GetProcAddress(module, "VT_GetDBSize_ENG");
    get_db_size_fn function;
    _Static_assert(sizeof(function) == sizeof(symbol), "PE32 function pointer size");
    memcpy(&function, &symbol, sizeof(function));
    return function;
}

static unload_ext_fn resolve_unload(HMODULE module) {
    FARPROC symbol = GetProcAddress(module, "VT_UNLOADTTS_EXT_ENG");
    unload_ext_fn function;
    _Static_assert(sizeof(function) == sizeof(symbol), "PE32 function pointer size");
    memcpy(&function, &symbol, sizeof(function));
    return function;
}

static unsigned int occupancy_mask(HMODULE module) {
    uintptr_t base = (uintptr_t)module;
    unsigned int mask = 0;
    for (unsigned int slot = 0; slot < 6; ++slot) {
        uintptr_t state = *(uintptr_t *)(base + 0xa0464U + slot * sizeof(uintptr_t));
        if (state != 0) {
            mask |= 1U << slot;
        }
    }
    return mask;
}

static void report_sizes(get_db_size_fn get_db_size, const char *label) {
    printf("%s", label);
    for (int slot = 0; slot < 6; ++slot) {
        int size = 0;
        uint32_t result = get_db_size(&size, slot);
        printf(" slot%d=%lu/%d", slot, (unsigned long)result, size);
    }
    putchar('\n');
}

int main(int argc, char **argv) {
    HMODULE module = LoadLibraryA("Z:\\samples\\vt_pau.dll");
    if (module == NULL) {
        printf("HARNESS load_library=%lu\n", (unsigned long)GetLastError());
        return 2;
    }
    load_ext_fn load_ext = resolve_load(module);
    get_db_size_fn get_db_size = resolve_size(module);
    unload_ext_fn unload_ext = resolve_unload(module);
    if (load_ext == NULL || get_db_size == NULL || unload_ext == NULL) {
        printf("HARNESS missing_export load=%d size=%d unload=%d\n", load_ext != NULL,
               get_db_size != NULL, unload_ext != NULL);
        FreeLibrary(module);
        return 3;
    }

    char database_base[] = "Z:\\work\\";
    if (argc == 2 && strcmp(argv[1], "reload") == 0) {
        short first = load_ext(NULL, 1, database_base, 0, UINT32_MAX, NULL, NULL, UINT32_MAX);
        printf("RELOAD first_ax=%d occupancy=%02x\n", (int)first, occupancy_mask(module));
        if (first != 0) {
            FreeLibrary(module);
            return 4;
        }
        char same_database_base[] = "Z:\\work\\";
        short duplicate = load_ext(NULL, 1, same_database_base, 0, UINT32_MAX, NULL, NULL,
                                   UINT32_MAX);
        printf("RELOAD duplicate_ax=%d occupancy=%02x\n", (int)duplicate,
               occupancy_mask(module));
        report_sizes(get_db_size, "RELOAD after_duplicate");
        unload_ext(1);
        printf("RELOAD after_unload occupancy=%02x\n", occupancy_mask(module));
        char reloaded_database_base[] = "Z:\\work\\";
        short reloaded = load_ext(NULL, 1, reloaded_database_base, 0, UINT32_MAX, NULL, NULL,
                                  UINT32_MAX);
        printf("RELOAD after_unload_ax=%d occupancy=%02x\n", (int)reloaded,
               occupancy_mask(module));
        report_sizes(get_db_size, "RELOAD after_reload");
        if (occupancy_mask(module) != 0) {
            unload_ext(1);
        }
    } else if (argc == 2 && strcmp(argv[1], "path-switch") == 0) {
        short first = load_ext(NULL, 1, database_base, 0, UINT32_MAX, NULL, NULL, UINT32_MAX);
        printf("PATH first_ax=%d occupancy=%02x\n", (int)first, occupancy_mask(module));
        if (first != 0) {
            FreeLibrary(module);
            return 7;
        }
        char different_database_base[] = "Z:\\work\\not-a-real-directory\\";
        short different = load_ext(NULL, 1, different_database_base, 0, UINT32_MAX, NULL, NULL,
                                   UINT32_MAX);
        printf("PATH different_base_ax=%d occupancy=%02x\n", (int)different,
               occupancy_mask(module));
        report_sizes(get_db_size, "PATH after_different_base");
        char original_database_base[] = "Z:\\work\\";
        short original = load_ext(NULL, 1, original_database_base, 0, UINT32_MAX, NULL, NULL,
                                  UINT32_MAX);
        printf("PATH original_base_again_ax=%d occupancy=%02x\n", (int)original,
               occupancy_mask(module));
        report_sizes(get_db_size, "PATH after_original_base");
        unload_ext(1);
    } else if (argc == 3 && strcmp(argv[1], "invalid-switch") == 0) {
        char *end = NULL;
        long input_slot = strtol(argv[2], &end, 10);
        if (end == argv[2] || *end != '\0' || input_slot < -100000 || input_slot > 100000) {
            puts("HARNESS invalid_slot_argument");
            FreeLibrary(module);
            return 11;
        }
        short first = load_ext(NULL, (int)input_slot, database_base, 0, UINT32_MAX, NULL, NULL,
                               UINT32_MAX);
        printf("INVALID_SWITCH input=%ld first_ax=%d occupancy=%02x\n", input_slot,
               (int)first, occupancy_mask(module));
        fflush(stdout);
        if (first != 0) {
            FreeLibrary(module);
            return 12;
        }
        char different_database_base[] = "Z:\\work\\not-a-real-directory\\";
        short different = load_ext(NULL, (int)input_slot, different_database_base, 0,
                                   UINT32_MAX, NULL, NULL, UINT32_MAX);
        printf("INVALID_SWITCH input=%ld different_base_ax=%d occupancy=%02x\n", input_slot,
               (int)different, occupancy_mask(module));
        fflush(stdout);
        report_sizes(get_db_size, "INVALID_SWITCH after_different_base");
        char recovery_database_base[] = "Z:\\work\\";
        short recovery = load_ext(NULL, 4, recovery_database_base, 0, UINT32_MAX, NULL, NULL,
                                  UINT32_MAX);
        printf("INVALID_SWITCH recovery_james_ax=%d occupancy=%02x\n", (int)recovery,
               occupancy_mask(module));
        report_sizes(get_db_size, "INVALID_SWITCH after_recovery_james");
        if (occupancy_mask(module) & (1U << 4)) {
            unload_ext(4);
        }
        unload_ext(1);
    } else if (argc == 2 && strcmp(argv[1], "multi") == 0) {
        char paul_base[] = "Z:\\work\\";
        short paul_first = load_ext(NULL, 1, paul_base, 0, UINT32_MAX, NULL, NULL, UINT32_MAX);
        printf("MULTI paul_first_ax=%d occupancy=%02x\n", (int)paul_first,
               occupancy_mask(module));
        if (paul_first != 0) {
            FreeLibrary(module);
            return 8;
        }
        char james_base[] = "Z:\\work\\";
        short james_second = load_ext(NULL, 4, james_base, 0, UINT32_MAX, NULL, NULL,
                                      UINT32_MAX);
        printf("MULTI james_second_ax=%d occupancy=%02x\n", (int)james_second,
               occupancy_mask(module));
        report_sizes(get_db_size, "MULTI after_both");
        if (james_second != 0) {
            unload_ext(1);
            FreeLibrary(module);
            return 9;
        }
        unload_ext(1);
        printf("MULTI after_unload_paul occupancy=%02x\n", occupancy_mask(module));
        report_sizes(get_db_size, "MULTI after_unload_paul");
        unload_ext(4);
        printf("MULTI after_unload_james occupancy=%02x\n", occupancy_mask(module));
        report_sizes(get_db_size, "MULTI after_unload_james");

        char james_first_base[] = "Z:\\work\\";
        short james_first = load_ext(NULL, 4, james_first_base, 0, UINT32_MAX, NULL, NULL,
                                    UINT32_MAX);
        printf("MULTI james_first_ax=%d occupancy=%02x\n", (int)james_first,
               occupancy_mask(module));
        if (james_first != 0) {
            FreeLibrary(module);
            return 10;
        }
        char paul_second_base[] = "Z:\\work\\";
        short paul_second = load_ext(NULL, 1, paul_second_base, 0, UINT32_MAX, NULL, NULL,
                                     UINT32_MAX);
        printf("MULTI paul_second_ax=%d occupancy=%02x\n", (int)paul_second,
               occupancy_mask(module));
        report_sizes(get_db_size, "MULTI reverse_after_both");
        unload_ext(4);
        printf("MULTI reverse_after_unload_james occupancy=%02x\n", occupancy_mask(module));
        report_sizes(get_db_size, "MULTI reverse_after_unload_james");
        unload_ext(1);
        printf("MULTI reverse_after_unload_paul occupancy=%02x\n", occupancy_mask(module));
        report_sizes(get_db_size, "MULTI reverse_after_unload_paul");
    } else if (argc == 2 && strcmp(argv[1], "multi-path-switch") == 0) {
        char paul_base[] = "Z:\\work\\";
        short paul = load_ext(NULL, 1, paul_base, 0, UINT32_MAX, NULL, NULL, UINT32_MAX);
        char james_base[] = "Z:\\work\\";
        short james = load_ext(NULL, 4, james_base, 0, UINT32_MAX, NULL, NULL, UINT32_MAX);
        printf("MULTI_PATH initial_paul_ax=%d initial_james_ax=%d occupancy=%02x\n",
               (int)paul, (int)james, occupancy_mask(module));
        if (paul != 0 || james != 0) {
            if (occupancy_mask(module) & (1U << 1)) unload_ext(1);
            if (occupancy_mask(module) & (1U << 4)) unload_ext(4);
            FreeLibrary(module);
            return 13;
        }
        char different_for_paul[] = "Z:\\work\\not-a-real-directory\\";
        short paul_switch = load_ext(NULL, 1, different_for_paul, 0, UINT32_MAX, NULL, NULL,
                                     UINT32_MAX);
        printf("MULTI_PATH paul_different_base_ax=%d occupancy=%02x\n", (int)paul_switch,
               occupancy_mask(module));
        report_sizes(get_db_size, "MULTI_PATH after_paul_different_base");
        char different_for_james[] = "Z:\\work\\not-a-real-directory\\";
        short james_switch = load_ext(NULL, 4, different_for_james, 0, UINT32_MAX, NULL, NULL,
                                      UINT32_MAX);
        printf("MULTI_PATH james_different_base_ax=%d occupancy=%02x\n", (int)james_switch,
               occupancy_mask(module));
        report_sizes(get_db_size, "MULTI_PATH after_james_different_base");
        char original_for_paul[] = "Z:\\work\\";
        short paul_original = load_ext(NULL, 1, original_for_paul, 0, UINT32_MAX, NULL, NULL,
                                       UINT32_MAX);
        char original_for_james[] = "Z:\\work\\";
        short james_original = load_ext(NULL, 4, original_for_james, 0, UINT32_MAX, NULL, NULL,
                                        UINT32_MAX);
        printf("MULTI_PATH original_paul_ax=%d original_james_ax=%d occupancy=%02x\n",
               (int)paul_original, (int)james_original, occupancy_mask(module));
        report_sizes(get_db_size, "MULTI_PATH after_original_bases");
        unload_ext(1);
        unload_ext(4);
    } else if (argc == 2) {
        char *end = NULL;
        long input_slot = strtol(argv[1], &end, 10);
        if (end == argv[1] || *end != '\0' || input_slot < -100000 || input_slot > 100000) {
            puts("HARNESS invalid_slot_argument");
            FreeLibrary(module);
            return 5;
        }
        short result = load_ext(NULL, (int)input_slot, database_base, 0, UINT32_MAX, NULL, NULL,
                                UINT32_MAX);
        printf("SLOT input=%ld load_ax=%d occupancy=%02x\n", input_slot, (int)result,
               occupancy_mask(module));
        report_sizes(get_db_size, "SLOT sizes");
        if (occupancy_mask(module) != 0) {
            unload_ext(1);
        }
    } else {
        puts("HARNESS usage=probe-load-ext-state.exe <slot|reload|path-switch|multi|multi-path-switch|invalid-switch slot>");
        FreeLibrary(module);
        return 6;
    }
    FreeLibrary(module);
    return 0;
}
