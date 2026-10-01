#define WIN32_LEAN_AND_MEAN
#include <windows.h>

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

static int replace_field(const char *input, unsigned int field_index, const char *replacement,
                         char *output, size_t capacity) {
    size_t input_length = strlen(input);
    if (input_length + 1 > capacity) return 0;
    memcpy(output, input, input_length + 1);
    char *start = output;
    for (unsigned int i = 0; i < field_index; ++i) {
        start = strchr(start, ':');
        if (start == NULL) return 0;
        ++start;
    }
    char *end = strchr(start, ':');
    if (end == NULL) return 0;
    size_t old_length = (size_t)(end - start);
    size_t replacement_length = strlen(replacement);
    if (old_length != replacement_length) return 0;
    memcpy(start, replacement, old_length);
    return 1;
}

static int replace_attribute(const char *input, const char *name, const char *replacement,
                             char *output, size_t capacity) {
    size_t input_length = strlen(input);
    if (input_length + 1 > capacity) return 0;
    memcpy(output, input, input_length + 1);
    char needle[48];
    if (snprintf(needle, sizeof(needle), "%s=\"", name) < 0) return 0;
    char *start = strstr(output, needle);
    if (start == NULL) return 0;
    start += strlen(needle);
    char *end = strchr(start, '"');
    if (end == NULL) return 0;
    size_t old_length = (size_t)(end - start);
    if (old_length != strlen(replacement)) return 0;
    memcpy(start, replacement, old_length);
    return 1;
}

static void run_case(check_license_fn check, const char *label, const char *input,
                     size_t input_length) {
    char *copy = malloc(input_length + 1);
    if (copy == NULL) {
        printf("LICENSE_FIELD case=%s harness_error=allocation\n", label);
        return;
    }
    memcpy(copy, input, input_length);
    copy[input_length] = '\0';
    int result = check(NULL, copy, (unsigned int)input_length, NULL);
    printf("LICENSE_FIELD case=%s result=%d\n", label, result);
    free(copy);
}

static void run_mutated_case(check_license_fn check, const char *label, const char *original,
                             size_t input_length, int (*mutate)(const char *, char *, size_t)) {
    char *copy = malloc(input_length + 1);
    if (copy == NULL) {
        printf("LICENSE_FIELD case=%s harness_error=allocation\n", label);
        return;
    }
    if (!mutate(original, copy, input_length + 1)) {
        printf("LICENSE_FIELD case=%s harness_error=mutation\n", label);
        free(copy);
        return;
    }
    int result = check(NULL, copy, (unsigned int)input_length, NULL);
    printf("LICENSE_FIELD case=%s result=%d\n", label, result);
    free(copy);
}

static char original[2048];

static int mutate_expiry_only(const char *input, char *output, size_t capacity) {
    return replace_field(input, 4, "1", output, capacity);
}
static int mutate_expdate_only(const char *input, char *output, size_t capacity) {
    return replace_attribute(input, "expdate", "1", output, capacity);
}
static int mutate_expiry_and_expdate(const char *input, char *output, size_t capacity) {
    return replace_field(input, 4, "1", output, capacity) &&
           replace_attribute(output, "expdate", "1", output, capacity);
}
static int mutate_vw_token(const char *input, char *output, size_t capacity) {
    return replace_field(input, 3, "VW_VTPAI", output, capacity);
}
static int mutate_channel_only(const char *input, char *output, size_t capacity) {
    return replace_field(input, 5, "7", output, capacity);
}
static int mutate_xml_channel_only(const char *input, char *output, size_t capacity) {
    return replace_attribute(input, "channel", "7", output, capacity);
}
static int mutate_channel_and_xml(const char *input, char *output, size_t capacity) {
    return replace_field(input, 5, "7", output, capacity) &&
           replace_attribute(output, "channel", "7", output, capacity);
}

int main(void) {
    HMODULE dll = LoadLibraryA("Z:\\samples\\vt_pau.dll");
    if (dll == NULL) {
        printf("HARNESS error=load_library code=%lu\n", (unsigned long)GetLastError());
        return 2;
    }
    check_license_fn check =
        (check_license_fn)GetProcAddress(dll, "VT_CheckLicense_ENG");
    if (check == NULL) {
        puts("HARNESS error=missing_export");
        FreeLibrary(dll);
        return 3;
    }
    size_t length = 0;
    if (!read_record(original, sizeof(original), &length)) {
        puts("HARNESS error=read_record");
        FreeLibrary(dll);
        return 4;
    }
    int baseline = check(NULL, original, (unsigned int)length, NULL);
    printf("LICENSE_FIELD baseline=%d bytes=%u\n", baseline, (unsigned int)length);
    run_case(check, "baseline_repeat", original, length);
    run_mutated_case(check, "expiry_position_only", original, length, mutate_expiry_only);
    run_mutated_case(check, "expdate_attribute_only", original, length, mutate_expdate_only);
    run_mutated_case(check, "expiry_position_and_expdate", original, length,
                     mutate_expiry_and_expdate);
    run_mutated_case(check, "vw_vtapi_token", original, length, mutate_vw_token);
    run_mutated_case(check, "channel_position_only", original, length, mutate_channel_only);
    run_mutated_case(check, "channel_attribute_only", original, length,
                     mutate_xml_channel_only);
    run_mutated_case(check, "channel_position_and_attribute", original, length,
                     mutate_channel_and_xml);
    FreeLibrary(dll);
    return 0;
}
