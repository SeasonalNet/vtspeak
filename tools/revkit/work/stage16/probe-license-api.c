#define WIN32_LEAN_AND_MEAN
#include <windows.h>

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef int(__cdecl *check_license_fn)(char *, char *, unsigned int, unsigned char *);
typedef int(__cdecl *get_license_info_fn)(int, char *, char *, unsigned int, unsigned char *, int);
typedef unsigned int(__cdecl *get_license_comment_fn)(char *, char *, int);
typedef int(__cdecl *get_tts_info_fn)(int, char *, void *, int);
typedef int(__cdecl *check_attribute_fn)(char *, unsigned char *, unsigned char *);
typedef short(__cdecl *load_tts_fn)(HWND, int, char *, char *);
typedef void(__cdecl *unload_tts_fn)(int);

static int read_record(char *buffer, size_t capacity, size_t *length) {
    FILE *file = fopen("Z:\\work\\data-common\\verify\\verification.txt", "rb");
    if (file == NULL) {
        return 0;
    }
    *length = fread(buffer, 1, capacity - 1, file);
    int ok = !ferror(file) && feof(file);
    fclose(file);
    if (!ok) {
        return 0;
    }
    while (*length > 0 && (buffer[*length - 1] == '\n' || buffer[*length - 1] == '\r')) {
        --*length;
    }
    buffer[*length] = '\0';
    return 1;
}

static int attr_value(const char *record, const char *name, char *out, size_t capacity) {
    char needle[48];
    if (snprintf(needle, sizeof(needle), "%s=\"", name) < 0) {
        return 0;
    }
    const char *start = strstr(record, needle);
    if (start == NULL) {
        return 0;
    }
    start += strlen(needle);
    const char *end = strchr(start, '"');
    if (end == NULL || (size_t)(end - start) + 1 > capacity) {
        return 0;
    }
    memcpy(out, start, (size_t)(end - start));
    out[end - start] = '\0';
    return 1;
}

static int field_value(const char *record, unsigned int wanted, char *out, size_t capacity) {
    const char *start = record;
    for (unsigned int field = 0; field < wanted; ++field) {
        start = strchr(start, ':');
        if (start == NULL) {
            return 0;
        }
        ++start;
    }
    const char *end = strchr(start, ':');
    if (end == NULL || (size_t)(end - start) + 1 > capacity) {
        return 0;
    }
    memcpy(out, start, (size_t)(end - start));
    out[end - start] = '\0';
    return 1;
}

static int mutate_attribute(const char *input, const char *name, const char *replacement,
                            char *output, size_t capacity) {
    size_t input_size = strlen(input);
    if (input_size + 1 > capacity) return 0;
    memcpy(output, input, input_size + 1);
    char needle[48];
    snprintf(needle, sizeof(needle), "%s=\"", name);
    char *start = strstr(output, needle);
    if (start == NULL) return 0;
    start += strlen(needle);
    char *end = strchr(start, '"');
    if (end == NULL) return 0;
    size_t old_size = (size_t)(end - start);
    size_t replacement_size = strlen(replacement);
    size_t tail_size = strlen(end) + 1;
    if (input_size - old_size + replacement_size + 1 > capacity) return 0;
    memmove(start + replacement_size, end, tail_size);
    memcpy(start, replacement, replacement_size);
    return 1;
}

static int remove_attribute(const char *input, const char *name, char *output, size_t capacity) {
    size_t input_size = strlen(input);
    if (input_size + 1 > capacity) return 0;
    memcpy(output, input, input_size + 1);
    char needle[48];
    snprintf(needle, sizeof(needle), "%s=\"", name);
    char *start = strstr(output, needle);
    if (start == NULL || start == output || start[-1] != ' ') return 0;
    --start;
    char *end = strchr(start + 1, '"');
    if (end == NULL || end[1] != ' ') return 0;
    end += 2;
    memmove(start, end, strlen(end) + 1);
    return 1;
}

static void check_case(check_license_fn check, const char *label, const char *value, size_t length) {
    char *copy = malloc(length + 1);
    if (copy == NULL) {
        printf("CHECK case=%s harness_error=allocation\n", label);
        return;
    }
    memcpy(copy, value, length);
    copy[length] = '\0';
    int result = check(NULL, copy, (unsigned int)length, NULL);
    printf("CHECK case=%s return=%d length=%u\n", label, result, (unsigned int)length);
    free(copy);
}

int main(void) {
    HMODULE dll = LoadLibraryA("Z:\\samples\\vt_pau.dll");
    if (dll == NULL) {
        printf("HARNESS error=load_library code=%lu\n", (unsigned long)GetLastError());
        return 2;
    }
    check_license_fn check = (check_license_fn)GetProcAddress(dll, "VT_CheckLicense_ENG");
    get_license_info_fn get_info = (get_license_info_fn)GetProcAddress(dll, "VT_GetLicenseInfo_ENG");
    get_license_comment_fn get_comment =
        (get_license_comment_fn)GetProcAddress(dll, "VT_GetLicenseComment_ENG");
    get_tts_info_fn get_tts_info = (get_tts_info_fn)GetProcAddress(dll, "VT_GetTTSInfo_ENG");
    if (check == NULL || get_info == NULL || get_comment == NULL || get_tts_info == NULL) {
        printf("HARNESS error=missing_export check=%d info=%d comment=%d tts_info=%d\n",
               check != NULL, get_info != NULL, get_comment != NULL, get_tts_info != NULL);
        FreeLibrary(dll);
        return 3;
    }

    char record[2048];
    size_t length = 0;
    if (!read_record(record, sizeof(record), &length)) {
        printf("HARNESS error=read_record\n");
        FreeLibrary(dll);
        return 4;
    }
    printf("HARNESS source_bytes=%u\n", (unsigned int)length);

    const char *source_path = "Z:\\work\\data-common\\verify\\verification.txt";
    int file_default = check(NULL, NULL, 0xffffffffU, NULL);
    int file_named = check("License", NULL, 0xffffffffU, NULL);
    int file_supplied = check((char *)source_path, NULL, 0xffffffffU, NULL);
    printf("CHECK file_default return=%d invalid_filename return=%d supplied_file return=%d\n",
           file_default, file_named, file_supplied);
    const char *info_sources[] = {"default", "supplied_file", "invalid_path"};
    char *info_paths[] = {NULL, (char *)source_path, "License"};
    for (size_t source = 0; source < sizeof(info_sources) / sizeof(info_sources[0]); ++source) {
        for (int request = 1; request <= 2; ++request) {
            uint32_t info_value = 0xffffffffU;
            int info_result = get_tts_info(request, info_paths[source], &info_value,
                                           sizeof(info_value));
            printf("TTSINFO request=%d source=%s result=%d value=%d\n", request,
                   info_sources[source], info_result, (int)info_value);
        }
    }
    char comment[2048] = {0};
    unsigned int comment_result = get_comment((char *)source_path, comment, sizeof(comment));
    char source_comment[1024];
    int comment_matches_source = field_value(record, 8, source_comment, sizeof(source_comment)) &&
                                strcmp(comment, source_comment) == 0;
    printf("COMMENT supplied_file result=%u output_bytes=%u matches_source_field=%d\n",
           comment_result, (unsigned int)strlen(comment), comment_matches_source);
    char host_id[128], xml_host_id[128], channel[32], xml_channel[32];
    char expiration[32], xml_expiration[32], user[128], xml_user[128];
    char platform[128], xml_os[128];
    int positions_match = field_value(record, 2, host_id, sizeof(host_id)) &&
                          attr_value(record, "hostid", xml_host_id, sizeof(xml_host_id)) &&
                          strcmp(host_id, xml_host_id) == 0 &&
                          field_value(record, 5, channel, sizeof(channel)) &&
                          attr_value(record, "channel", xml_channel, sizeof(xml_channel)) &&
                          strcmp(channel, xml_channel) == 0 &&
                          field_value(record, 4, expiration, sizeof(expiration)) &&
                          attr_value(record, "expdate", xml_expiration, sizeof(xml_expiration)) &&
                          strcmp(expiration, xml_expiration) == 0 &&
                          field_value(record, 6, user, sizeof(user)) &&
                          attr_value(record, "user", xml_user, sizeof(xml_user)) &&
                          strcmp(user, xml_user) == 0 &&
                          field_value(record, 7, platform, sizeof(platform)) &&
                          attr_value(record, "os", xml_os, sizeof(xml_os)) &&
                          strcmp(platform, xml_os) == 0;
    printf("POSITION_SAMPLE hostid_channel_expdate_user_os_crossmatch=%d\n", positions_match);

    check_attribute_fn check_attribute =
        (check_attribute_fn)((uintptr_t)dll + 0x29d30U);
    char tag_attributes[1024];
    const char *tag_start = strstr(record, "<vw_verify ");
    const char *tag_end = tag_start == NULL ? NULL : strstr(tag_start, "/>" );
    if (tag_start != NULL && tag_end != NULL) {
        tag_start += strlen("<vw_verify ");
        size_t tag_size = (size_t)(tag_end - tag_start);
        if (tag_size < sizeof(tag_attributes)) {
            memcpy(tag_attributes, tag_start, tag_size);
            tag_attributes[tag_size] = '\0';
            const char *attr_keys[] = {"os", "lang", "speaker", "version", "dbaccess",
                                       "sampling", "dbsize"};
            const char *attr_bad[] = {"linux", "unknown", "unknown", "0", "none", "8", "0"};
            char speaker_variant[1024];
            for (size_t i = 0; i < sizeof(attr_keys) / sizeof(attr_keys[0]); ++i) {
                int baseline = check_attribute(tag_attributes, (unsigned char *)attr_keys[i], NULL);
                char variant[1024];
                int changed = mutate_attribute(tag_attributes, attr_keys[i], attr_bad[i],
                                               variant, sizeof(variant));
                int altered = changed ? check_attribute(variant, (unsigned char *)attr_keys[i], NULL)
                                      : -1;
                int removed = remove_attribute(tag_attributes, attr_keys[i], variant,
                                               sizeof(variant));
                int absent = removed ? check_attribute(variant, (unsigned char *)attr_keys[i], NULL)
                                     : -1;
                printf("ATTR_HELPER key=%s baseline=%d altered=%d missing=%d\n", attr_keys[i],
                       baseline, altered, absent);
            }
            int speaker_known = mutate_attribute(tag_attributes, "speaker", "paul", speaker_variant,
                                                 sizeof(speaker_variant));
            int speaker_known_result = speaker_known
                                           ? check_attribute(speaker_variant,
                                                             (unsigned char *)"speaker", NULL)
                                           : -1;
            int speaker_unknown = mutate_attribute(tag_attributes, "speaker", "not-a-voice",
                                                   speaker_variant, sizeof(speaker_variant));
            int speaker_unknown_result = speaker_unknown
                                             ? check_attribute(speaker_variant,
                                                               (unsigned char *)"speaker", NULL)
                                             : -1;
            printf("ATTR_HELPER key=speaker known_value=%d unknown_value=%d\n",
                   speaker_known_result, speaker_unknown_result);
            char dbsize_unknown[1024];
            int dbsize_mutated = mutate_attribute(tag_attributes, "dbsize", "0", dbsize_unknown,
                                                 sizeof(dbsize_unknown));
            char dbsize_large[1024];
            int dbsize_large_mutated = mutate_attribute(tag_attributes, "dbsize", "999999999",
                                                        dbsize_large, sizeof(dbsize_large));
            const char *allowed_speakers[] = {"kate", "paul", "julie", "james", "violeta",
                                              "bridget"};
            for (unsigned int slot = 0; slot < 6; ++slot) {
                const char *voice_constraint =
                    *(const char **)((uintptr_t)dll + 0x7c6a8U + slot * 6U * sizeof(void *));
                int list_match = 0;
                for (size_t allowed = 0;
                     allowed < sizeof(allowed_speakers) / sizeof(allowed_speakers[0]); ++allowed) {
                    if (_stricmp(voice_constraint, allowed_speakers[allowed]) == 0) {
                        list_match = 1;
                        break;
                    }
                }
                int speaker_result = check_attribute(tag_attributes, (unsigned char *)"speaker",
                                                     (unsigned char *)voice_constraint);
                int dbsize_result =
                    check_attribute(tag_attributes, (unsigned char *)"dbsize",
                                    (unsigned char *)voice_constraint);
                int dbsize_zero_result =
                    dbsize_mutated ? check_attribute(dbsize_unknown, (unsigned char *)"dbsize",
                                                     (unsigned char *)voice_constraint)
                                   : -1;
                int dbsize_large_result =
                    dbsize_large_mutated ? check_attribute(dbsize_large, (unsigned char *)"dbsize",
                                                          (unsigned char *)voice_constraint)
                                         : -1;
                printf("ATTR_CONTEXT slot=%u speaker=%d speaker_list_match=%d dbsize=%d "
                       "dbsize_zero=%d dbsize_large=%d\n", slot, speaker_result, list_match,
                       dbsize_result, dbsize_zero_result, dbsize_large_result);
            }
        } else {
            printf("ATTR_HELPER harness_error=tag_too_long\n");
        }
    } else {
        printf("ATTR_HELPER harness_error=tag_boundary\n");
    }
    check_case(check, "supplied_record", record, length);

    char expected[1024];
    char output[1024];
    const struct {
        int selector;
        const char *name;
    } selectors[] = {
        {0, "channel_number"}, {1, "expiry_token"}, {2, "host_id"},
        {3, "tag_attributes"}, {4, "os"},           {5, "lang"},
        {6, "speaker"},        {7, "version"},      {8, "dbaccess"},
        {9, "sampling"},       {10, "app"},         {11, "wavsave"},
        {12, "savetime"},      {13, "bgaudio"},     {14, "dbsize"},
        {15, "realtime"},
    };
    const char *attribute_names[] = {"", "", "", "", "os", "lang", "speaker", "version",
                                     "dbaccess", "sampling", "app", "wavsave", "", "bgaudio",
                                     "dbsize", "realtime"};
    for (size_t i = 0; i < sizeof(selectors) / sizeof(selectors[0]); ++i) {
        memset(output, 0, sizeof(output));
        int result = get_info(selectors[i].selector, NULL, record, (unsigned int)length,
                              (unsigned char *)output, sizeof(output));
        int matches = -1;
        if (selectors[i].selector == 0 || selectors[i].selector == 12) {
            uint32_t got = 0;
            memcpy(&got, output, sizeof(got));
            if (selectors[i].selector == 0 && field_value(record, 5, expected, sizeof(expected))) {
                matches = got == (uint32_t)strtoul(expected, NULL, 10);
            } else if (selectors[i].selector == 12 && attr_value(record, "savetime", expected, sizeof(expected))) {
                matches = got == (uint32_t)strtoul(expected, NULL, 10);
            } else {
                matches = 0;
            }
        } else if (selectors[i].selector == 1) {
            matches = field_value(record, 4, expected, sizeof(expected)) && strcmp(output, expected) == 0;
        } else if (selectors[i].selector == 2) {
            matches = field_value(record, 2, expected, sizeof(expected)) && strcmp(output, expected) == 0;
        } else if (selectors[i].selector == 3) {
            const char *tag = strstr(record, "<vw_verify ");
            const char *end = tag == NULL ? NULL : strstr(tag, "/>" );
            int full_match = 0;
            int body_match = 0;
            int attributes_match = 0;
            if (tag != NULL && end != NULL && (size_t)(end - tag) + 1 < sizeof(expected)) {
                size_t tag_length = (size_t)(end - tag) + 2;
                memcpy(expected, tag, tag_length);
                expected[tag_length] = '\0';
                full_match = strcmp(output, expected) == 0;
                const char *body = tag + strlen("<vw_verify ");
                size_t body_length = (size_t)(end - body) + 2;
                memcpy(expected, body, body_length);
                expected[body_length] = '\0';
                body_match = strcmp(output, expected) == 0;
                size_t attr_length = (size_t)(end - body);
                memcpy(expected, body, attr_length);
                expected[attr_length] = '\0';
                attributes_match = strcmp(output, expected) == 0;
                matches = full_match || body_match || attributes_match;
            } else {
                matches = 0;
            }
            printf("INFO_TAG full_match=%d body_with_close_match=%d attribute_text_match=%d\n",
                   full_match, body_match, attributes_match);
        } else {
            matches = attr_value(record, attribute_names[selectors[i].selector], expected,
                                 sizeof(expected))
                          ? strcmp(output, expected) == 0
                          : output[0] == 0;
        }
        size_t output_size = selectors[i].selector == 0 || selectors[i].selector == 12
                                 ? sizeof(uint32_t)
                                 : strlen(output) + 1;
        printf("INFO selector=%d name=%s return=%d output_bytes=%u source_match=%d\n",
               selectors[i].selector, selectors[i].name, result, (unsigned int)output_size, matches);
        memset(expected, 0, sizeof(expected));
        int file_result = get_info(selectors[i].selector, (char *)source_path, NULL, 0xffffffffU,
                                   (unsigned char *)expected, sizeof(expected));
        int file_match = result == file_result &&
                         memcmp(output, expected, output_size) == 0;
        printf("INFO_FILE selector=%d return=%d matches_memory=%d\n", selectors[i].selector,
               file_result, file_match);

        for (int file_mode = 0; file_mode <= 1; ++file_mode) {
            int first_success = -1;
            int undersized_result = 0;
            int at_321_result = 0;
            int guards_ok = 1;
            int capacity_limit = 512;
            int final_result = 0;
            for (int capacity = 0; capacity <= capacity_limit; ++capacity) {
                unsigned char guarded[1024];
                memset(guarded, 0xa5, sizeof(guarded));
                int capacity_result = file_mode
                                          ? get_info(selectors[i].selector, (char *)source_path,
                                                     NULL, 0xffffffffU, guarded, capacity)
                                          : get_info(selectors[i].selector, NULL, record,
                                                     (unsigned int)length, guarded, capacity);
                for (size_t guard = (size_t)capacity; guard < sizeof(guarded); ++guard) {
                    if (guarded[guard] != 0xa5) guards_ok = 0;
                }
                if (capacity_result == 0 && first_success < 0) first_success = capacity;
                if (capacity == (int)output_size - 1) undersized_result = capacity_result;
                if (capacity == 321) at_321_result = capacity_result;
                if (capacity == capacity_limit) final_result = capacity_result;
            }
            printf("INFO_CAP_MATRIX selector=%d source=%s output_bytes=%u first_success=%d undersized=%d at_321=%d final=%d guards_ok=%d\n",
                   selectors[i].selector, file_mode ? "file" : "memory",
                   (unsigned int)output_size, first_success, undersized_result, at_321_result,
                   final_result,
                   guards_ok);
            int size_query = get_info(selectors[i].selector,
                                      file_mode ? (char *)source_path : NULL,
                                      file_mode ? NULL : record,
                                      file_mode ? 0xffffffffU : (unsigned int)length,
                                      NULL, -1);
            int null_output = get_info(selectors[i].selector,
                                       file_mode ? (char *)source_path : NULL,
                                       file_mode ? NULL : record,
                                       file_mode ? 0xffffffffU : (unsigned int)length,
                                       NULL, 512);
            printf("INFO_EDGE_MATRIX selector=%d source=%s size_query=%d null_output=%d\n",
                   selectors[i].selector, file_mode ? "file" : "memory", size_query,
                   null_output);
        }
    }

    const int cap_selector[] = {0, 0, 1, 1, 1, 99};
    const int cap_size[] = {3, 4, 0, 1, 2, 16};
    for (size_t i = 0; i < sizeof(cap_selector) / sizeof(cap_selector[0]); ++i) {
        int selector = cap_selector[i];
        int capacity = cap_size[i];
        unsigned char tiny[16] = {0};
        int result = get_info(selector, NULL, record, (unsigned int)length, tiny, capacity);
        printf("INFO_CAP selector=%d capacity=%d return=%d\n", selector, capacity, result);
    }
    int null_output = get_info(1, NULL, record, (unsigned int)length, NULL, 16);
    int invalid_selector = get_info(16, NULL, record, (unsigned int)length, (unsigned char *)output,
                                    sizeof(output));
    printf("INFO_EDGE null_output=%d invalid_selector=%d\n", null_output, invalid_selector);

    for (int file_mode = 0; file_mode <= 1; ++file_mode) {
        int first_success = -1;
        int below_result = 0;
        int guards_ok = 1;
        for (int capacity = 0; capacity <= (int)sizeof(comment); ++capacity) {
            unsigned char guarded[2048];
            memset(guarded, 0xa5, sizeof(guarded));
            unsigned int result = get_comment(file_mode ? (char *)source_path : NULL,
                                              (char *)guarded, capacity);
            if ((int)result >= 0 && first_success < 0) first_success = capacity;
            if (capacity == 321) below_result = (int)result;
            for (size_t guard = (size_t)capacity; guard < sizeof(guarded); ++guard) {
                if (guarded[guard] != 0xa5) guards_ok = 0;
            }
        }
        printf("COMMENT_CAP_MATRIX source=%s first_success=%d below_full=%d guards_ok=%d\n",
               file_mode ? "file" : "default", first_success, below_result, guards_ok);
    }

    const char *field_names[] = {"digest_changed", "host_id_changed", "product_changed",
                                 "field_4_changed", "field_5_changed", "user_changed",
                                 "platform_changed"};
    const unsigned int field_positions[] = {1, 2, 3, 4, 5, 6, 7};
    for (size_t i = 0; i < sizeof(field_names) / sizeof(field_names[0]); ++i) {
        char mutation[2048];
        memcpy(mutation, record, length + 1);
        char *start = mutation;
        for (unsigned int field = 0; field < field_positions[i]; ++field) {
            start = strchr(start, ':');
            if (start == NULL) break;
            ++start;
        }
        if (start != NULL) {
            char *end = strchr(start, ':');
            if (end != NULL && end > start) {
                start[0] = start[0] == '0' ? '1' : '0';
                check_case(check, field_names[i], mutation, length);
            }
        }
    }

    const char *attrs[] = {"app", "expdate", "os", "speaker", "version", "dbaccess", "sampling", "dbsize"};
    const char *replacements[] = {"unknown", "20000101", "unknown", "unknown", "0", "unknown", "8", "0"};
    for (size_t i = 0; i < sizeof(attrs) / sizeof(attrs[0]); ++i) {
        char mutation[2048];
        memcpy(mutation, record, length + 1);
        char pattern[48];
        snprintf(pattern, sizeof(pattern), "%s=\"", attrs[i]);
        char *start = strstr(mutation, pattern);
        if (start == NULL) continue;
        start += strlen(pattern);
        char *end = strchr(start, '"');
        if (end == NULL) continue;
        size_t old_size = (size_t)(end - start);
        size_t replacement_size = strlen(replacements[i]);
        if (replacement_size > old_size) continue;
        memcpy(start, replacements[i], replacement_size);
        memset(start + replacement_size, ' ', old_size - replacement_size);
        char label[64];
        snprintf(label, sizeof(label), "attr_%s_changed", attrs[i]);
        check_case(check, label, mutation, length);
    }

    const char *remove_attrs[] = {"os", "speaker", "version", "dbaccess", "sampling", "dbsize"};
    for (size_t i = 0; i < sizeof(remove_attrs) / sizeof(remove_attrs[0]); ++i) {
        char mutation[2048];
        memcpy(mutation, record, length + 1);
        char pattern[48];
        snprintf(pattern, sizeof(pattern), " %s=\"", remove_attrs[i]);
        char *start = strstr(mutation, pattern);
        if (start == NULL) continue;
        char *end = strchr(start + strlen(pattern), '"');
        if (end == NULL || end[1] != ' ') continue;
        end += 2;
        memmove(start, end, strlen(end) + 1);
        size_t mutated_length = length - (size_t)(end - start);
        char label[64];
        snprintf(label, sizeof(label), "attr_%s_missing", remove_attrs[i]);
        check_case(check, label, mutation, mutated_length);
    }

    char malformed[2048];
    memcpy(malformed, record, length + 1);
    char *tag_open = strstr(malformed, "<vw_verify ");
    if (tag_open != NULL) tag_open[1] = 'x';
    check_case(check, "tag_open_changed", malformed, length);
    memcpy(malformed, record, length + 1);
    char *tag_close = strstr(malformed, "/>" );
    if (tag_close != NULL) tag_close[1] = ' ';
    check_case(check, "tag_close_changed", malformed, length);
    check_case(check, "short_input", "short", 5);

    char framed[4096];
    framed[0] = 'X'; framed[1] = 'Y'; framed[2] = 'Z';
    for (size_t i = 0; i < length; ++i) framed[3 + i] = (char)((unsigned char)record[i] + 1U);
    framed[3 + length] = 1;
    framed[4 + length] = 'X'; framed[5 + length] = 'Y'; framed[6 + length] = 'Z';
    check_case(check, "framed_valid", framed, length + 7);
    const char *framed_path = "Z:\\tmp\\vtspeak-license-framed.tmp";
    FILE *framed_file = fopen(framed_path, "wb");
    int framed_file_written = framed_file != NULL &&
                              fwrite(framed, 1, length + 7, framed_file) == length + 7;
    if (framed_file != NULL) fclose(framed_file);
    if (framed_file_written) {
        int framed_file_result = check((char *)framed_path, NULL, 0xffffffffU, NULL);
        char file_output[1024] = {0};
        int framed_info_result = get_info(2, (char *)framed_path, NULL, 0xffffffffU,
                                          (unsigned char *)file_output, sizeof(file_output));
        unsigned int framed_comment_result = get_comment((char *)framed_path, comment,
                                                         sizeof(comment));
        printf("CHECK framed_file return=%d info_selector_2=%d host_match=%d comment_match=%d\n",
               framed_file_result, framed_info_result,
               strcmp(file_output, host_id) == 0,
               strcmp(comment, source_comment) == 0 && framed_comment_result == comment_result);
        DeleteFileA(framed_path);
    } else {
        if (framed_file != NULL) DeleteFileA(framed_path);
        printf("CHECK framed_file harness_error=write\n");
    }
    framed[3 + length] = 2;
    check_case(check, "framed_wrong_seed", framed, length + 7);
    framed[3 + length] = 1;
    framed[6 + length] = 'Q';
    check_case(check, "framed_bad_suffix", framed, length + 7);

    FARPROC load_address = GetProcAddress(dll, "VT_LOADTTS_ENG");
    FARPROC unload_address = GetProcAddress(dll, "VT_UNLOADTTS_ENG");
    load_tts_fn load_tts = NULL;
    unload_tts_fn unload_tts = NULL;
    memcpy(&load_tts, &load_address, sizeof(load_tts));
    memcpy(&unload_tts, &unload_address, sizeof(unload_tts));
    if (load_address != NULL && unload_address != NULL) {
        short load_result = load_tts(NULL, -1, NULL, NULL);
        printf("LOAD default_slot result=%d\n", load_result);
        if (load_result == 0) {
            for (unsigned int slot = 0; slot < 6; ++slot) {
                const char *voice_constraint =
                    *(const char **)((uintptr_t)dll + 0x7c6a8U + slot * 6U * sizeof(void *));
                int speaker_result = check_attribute(tag_attributes, (unsigned char *)"speaker",
                                                     (unsigned char *)voice_constraint);
                int dbsize_result = check_attribute(tag_attributes, (unsigned char *)"dbsize",
                                                    (unsigned char *)voice_constraint);
                printf("ATTR_LOADED slot=%u speaker=%d dbsize=%d\n", slot, speaker_result,
                       dbsize_result);
                if (slot == 1) {
                    int loaded_file_check = check((char *)source_path, NULL, 0xffffffffU,
                                                  (unsigned char *)voice_constraint);
                    printf("CHECK loaded_paul_file return=%d\n", loaded_file_check);
                    const char *limits[] = {"300", "450", "451", "452", "500"};
                    for (size_t limit = 0; limit < sizeof(limits) / sizeof(limits[0]); ++limit) {
                        char limited_attributes[1024];
                        if (mutate_attribute(tag_attributes, "dbsize", limits[limit],
                                             limited_attributes, sizeof(limited_attributes))) {
                            int allowed = check_attribute(limited_attributes,
                                                          (unsigned char *)"dbsize",
                                                          (unsigned char *)voice_constraint);
                            printf("ATTR_DB_LIMIT value=%s accepted=%d\n", limits[limit], allowed);
                        }
                    }
                }
            }
            unload_tts(-1);
            printf("UNLOAD default_slot completed=1\n");
        }
    } else {
        printf("LOAD harness_error=missing_export\n");
    }

    FreeLibrary(dll);
    return 0;
}
