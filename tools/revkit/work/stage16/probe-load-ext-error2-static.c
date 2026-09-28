#include <stdint.h>
#include <stdio.h>

#include <windows.h>

__declspec(dllimport) short __cdecl VT_LOADTTS_EXT_ENG(
    HWND window, int slot, char *database_base, uint32_t argument4, uint32_t argument5,
    char *license_path, char *license_text, uint32_t license_length);

int main(void) {
    char database_base[] = "Z:\\work\\";
    short result = VT_LOADTTS_EXT_ENG(NULL, 1, database_base, 0, UINT32_MAX, NULL, NULL,
                                      UINT32_MAX);
    printf("ERROR2_PROBE exported_load_ax=%d\n", (int)result);
    fflush(stdout);
    return 0;
}
