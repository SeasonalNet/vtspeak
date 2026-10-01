#define WIN32_LEAN_AND_MEAN
#include <windows.h>

#include <ctype.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef int(__cdecl *check_license_fn)(char *, char *, unsigned int, unsigned char *);

static int read_record(char *buffer, size_t capacity, size_t *length) {
    FILE *file = fopen("Z:\\work\\data-common\\verify\\verification.txt", "rb");
    if (file == NULL) return 0;
    *length = fread(buffer, 1, capacity - 1, file);
    int ok = !ferror(file) && feof(file);
    fclose(file);
    if (!ok) return 0;
    while (*length > 0 && (buffer[*length - 1] == '\n' || buffer[*length - 1] == '\r')) {
        --*length;
    }
    buffer[*length] = '\0';
    return 1;
}

static int is_hex_field(const char *start, size_t length) {
    if (length != 96) return 0;
    for (size_t i = 0; i < length; ++i) {
        if (!isxdigit((unsigned char)start[i])) return 0;
    }
    return 1;
}

static size_t find_token(const char *record, size_t length, size_t *matches) {
    size_t found = 0;
    size_t found_at = 0;
    for (size_t start = 0; start + 96 <= length; ++start) {
        if (!is_hex_field(record + start, 96)) continue;
        int left_ok = start == 0 || record[start - 1] == ':' || record[start - 1] == ';';
        int right_ok = start + 96 == length || record[start + 96] == ':' ||
                       record[start + 96] == ';';
        if (!left_ok || !right_ok) continue;
        found_at = start;
        ++found;
    }
    *matches = found;
    return found_at;
}

static char next_hex_digit(char value) {
    static const char digits[] = "0123456789abcdef";
    const char *at = strchr(digits, (char)tolower((unsigned char)value));
    if (at == NULL) return '0';
    size_t index = (size_t)(at - digits);
    return digits[(index + 1) & 15U];
}

static int flip_hex_case(char *value) {
    if (*value >= 'a' && *value <= 'f') {
        *value = (char)(*value - 'a' + 'A');
        return 1;
    }
    if (*value >= 'A' && *value <= 'F') {
        *value = (char)(*value - 'A' + 'a');
        return 1;
    }
    return 0;
}

static int selected_position(size_t token_offset) {
    return token_offset % 6U < 2U;
}

int main(void) {
    HMODULE dll = LoadLibraryA("Z:\\samples\\vt_pau.dll");
    if (dll == NULL) {
        printf("HARNESS error=load_library code=%lu\n", (unsigned long)GetLastError());
        return 2;
    }
    check_license_fn check = (check_license_fn)GetProcAddress(dll, "VT_CheckLicense_ENG");
    if (check == NULL) {
        puts("HARNESS error=missing_export");
        FreeLibrary(dll);
        return 3;
    }

    char original[2048];
    size_t length = 0;
    if (!read_record(original, sizeof(original), &length)) {
        puts("HARNESS error=read_record");
        FreeLibrary(dll);
        return 4;
    }
    size_t matches = 0;
    size_t token_start = find_token(original, length, &matches);
    if (matches != 1) {
        printf("HARNESS error=token_count count=%u\n", (unsigned int)matches);
        FreeLibrary(dll);
        return 5;
    }

    int baseline = check(NULL, original, (unsigned int)length, NULL);
    printf("TOKEN_MATRIX baseline=%d record_bytes=%u token_length=96\n", baseline,
           (unsigned int)length);

    unsigned int hex_selected_pass = 0, hex_selected_fail = 0;
    unsigned int hex_gap_pass = 0, hex_gap_fail = 0;
    unsigned int nonhex_selected_pass = 0, nonhex_selected_fail = 0;
    unsigned int nonhex_gap_pass = 0, nonhex_gap_fail = 0;
    unsigned int caseflip_count = 0, caseflip_pass = 0, caseflip_fail = 0;
    size_t accepted_flip_positions[96];
    size_t accepted_flip_count = 0;
    for (size_t index = 0; index < 96; ++index) {
        int selected = selected_position(index);
        char mutated[2048];
        memcpy(mutated, original, length + 1);
        mutated[token_start + index] = next_hex_digit(mutated[token_start + index]);
        int result = check(NULL, mutated, (unsigned int)length, NULL);
        if (selected) {
            if (result == baseline) ++hex_selected_pass;
            else ++hex_selected_fail;
        } else {
            if (result == baseline) ++hex_gap_pass;
            else ++hex_gap_fail;
        }
        printf("TOKEN_BYTE index=%02u class=%s mutation=next_hex result=%d\n",
               (unsigned int)index, selected ? "selected" : "gap", result);

        memcpy(mutated, original, length + 1);
        mutated[token_start + index] = 'g';
        result = check(NULL, mutated, (unsigned int)length, NULL);
        if (selected) {
            if (result == baseline) ++nonhex_selected_pass;
            else ++nonhex_selected_fail;
        } else {
            if (result == baseline) ++nonhex_gap_pass;
            else ++nonhex_gap_fail;
        }
        printf("TOKEN_BYTE index=%02u class=%s mutation=nonhex result=%d\n",
               (unsigned int)index, selected ? "selected" : "gap", result);

        memcpy(mutated, original, length + 1);
        if (flip_hex_case(&mutated[token_start + index])) {
            ++caseflip_count;
            result = check(NULL, mutated, (unsigned int)length, NULL);
            if (result == baseline) {
                ++caseflip_pass;
                accepted_flip_positions[accepted_flip_count++] = index;
                printf("TOKEN_CASEFLIP_ACCEPTED index=%02u\n", (unsigned int)index);
            } else {
                ++caseflip_fail;
            }
            printf("TOKEN_BYTE index=%02u class=%s mutation=case_flip result=%d\n",
                   (unsigned int)index, selected ? "selected" : "gap", result);
        }
    }

    unsigned int combo_count = 0, combo_pass = 0, combo_fail = 0;
    if (accepted_flip_count <= 16) {
        unsigned int limit = 1U << accepted_flip_count;
        for (unsigned int mask = 0; mask < limit; ++mask) {
            char mutated[2048];
            memcpy(mutated, original, length + 1);
            for (size_t bit = 0; bit < accepted_flip_count; ++bit) {
                if ((mask & (1U << bit)) != 0) {
                    (void)flip_hex_case(&mutated[token_start + accepted_flip_positions[bit]]);
                }
            }
            int result = check(NULL, mutated, (unsigned int)length, NULL);
            ++combo_count;
            if (result == baseline) ++combo_pass;
            else ++combo_fail;
            printf("TOKEN_CASE_COMBO mask=%u result=%d\n", mask, result);
        }
    }
    printf("TOKEN_MATRIX_SUMMARY selected_hex_pass=%u selected_hex_fail=%u "
           "gap_hex_pass=%u gap_hex_fail=%u selected_nonhex_pass=%u "
           "selected_nonhex_fail=%u gap_nonhex_pass=%u gap_nonhex_fail=%u "
           "caseflip_count=%u caseflip_pass=%u caseflip_fail=%u "
           "combo_count=%u combo_pass=%u combo_fail=%u\n",
           hex_selected_pass, hex_selected_fail, hex_gap_pass, hex_gap_fail,
           nonhex_selected_pass, nonhex_selected_fail, nonhex_gap_pass, nonhex_gap_fail,
           caseflip_count, caseflip_pass, caseflip_fail, combo_count, combo_pass, combo_fail);
    FreeLibrary(dll);
    return 0;
}
