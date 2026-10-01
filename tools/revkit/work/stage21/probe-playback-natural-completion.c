#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <mmsystem.h>

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef short(__cdecl *load_tts_fn)(HWND, int, char *, char *);
typedef short(__cdecl *play_tts_fn)(HWND, UINT, char *, int, int, int, int, int, int,
                                    int);
typedef void(__cdecl *stop_tts_fn)(void);
typedef void(__cdecl *unload_tts_fn)(int);
typedef int(__cdecl *get_info_fn)(int, char *, void *, int);

static void *resolve(HMODULE dll, const char *name, size_t target_size,
                     void *target_storage) {
    FARPROC address = GetProcAddress(dll, name);
    if (address == NULL) {
        printf("RESOLVE name=%s error=%lu\n", name, GetLastError());
        return NULL;
    }
    if (target_size != sizeof(address)) {
        printf("RESOLVE name=%s error=pointer_size_mismatch\n", name);
        return NULL;
    }
    memcpy(target_storage, &address, target_size);
    return target_storage;
}

static int query_state(get_info_fn get_info, const char *phase) {
    DWORD output = 0x5a5a5a5aU;
    int result = get_info(101, NULL, &output, (int)sizeof(output));
    printf("STATE phase=%s result=%d output=0x%08lx tick=%lu\n", phase,
           result, (unsigned long)output, (unsigned long)GetTickCount());
    return result;
}

int main(int argc, char **argv) {
    const char *dll_path = "Z:\\samples\\vt_pau.dll";
    char database_path[] = "Z:\\work\\";
    char text[1024];
    const char *source_text = argc > 1 ? argv[1] : "Hello there.";
    const UINT caller_message = WM_APP + 0x51U;
    HMODULE dll = LoadLibraryA(dll_path);
    load_tts_fn load_tts = NULL;
    play_tts_fn play_tts = NULL;
    stop_tts_fn stop_tts = NULL;
    unload_tts_fn unload_tts = NULL;
    get_info_fn get_info = NULL;
    int loaded = 0;
    int played = 0;
    int open_messages = 0;
    int done_messages = 0;
    int caller_messages = 0;
    int quit_seen = 0;

    if (strncmp(source_text, "hex:", 4) == 0) {
        const char *hex = source_text + 4;
        size_t hex_length = strlen(hex);
        if ((hex_length & 1U) != 0 || hex_length / 2U >= sizeof(text)) {
            printf("PROBE_INPUT_ERROR hex_length=%lu\n",
                   (unsigned long)hex_length);
            return 1;
        }
        for (size_t i = 0; i < hex_length / 2U; ++i) {
            char pair[3] = {hex[i * 2U], hex[i * 2U + 1U], '\0'};
            unsigned long value = strtoul(pair, NULL, 16);
            if (value == 0 || value > 255U) {
                printf("PROBE_INPUT_ERROR byte_index=%lu value=%lu\n",
                       (unsigned long)i, value);
                return 1;
            }
            text[i] = (char)value;
        }
        text[hex_length / 2U] = '\0';
    } else {
        snprintf(text, sizeof(text), "%s", source_text);
    }
    printf("PROBE_BEGIN dll=%s db=%s text=%s bytes=", dll_path, database_path,
           text);
    for (size_t i = 0; text[i] != '\0'; ++i) {
        printf("%02x", (unsigned char)text[i]);
    }
    printf("\n");
    if (dll == NULL) {
        printf("LOAD_LIBRARY error=%lu\n", GetLastError());
        return 2;
    }
    if (resolve(dll, "VT_LOADTTS_ENG", sizeof(load_tts), &load_tts) == NULL ||
        resolve(dll, "VT_PLAYTTS_ENG", sizeof(play_tts), &play_tts) == NULL ||
        resolve(dll, "VT_STOPTTS_ENG", sizeof(stop_tts), &stop_tts) == NULL ||
        resolve(dll, "VT_UNLOADTTS_ENG", sizeof(unload_tts), &unload_tts) == NULL ||
        resolve(dll, "VT_GetTTSInfo_ENG", sizeof(get_info), &get_info) == NULL) {
        FreeLibrary(dll);
        return 3;
    }

    short load_result = load_tts(NULL, 1, database_path, NULL);
    printf("LOAD result=%d\n", load_result);
    if (load_result != 0) {
        FreeLibrary(dll);
        return 4;
    }
    loaded = 1;
    query_state(get_info, "loaded_idle");

    short play_result = play_tts(NULL, caller_message, text, 1, 100, 100, 100,
                                 0, -1, 0);
    printf("PLAY result=%d\n", play_result);
    if (play_result != 1) {
        unload_tts(1);
        FreeLibrary(dll);
        return 5;
    }
    played = 1;
    query_state(get_info, "play_return");

    DWORD start_tick = GetTickCount();
    DWORD deadline = start_tick + 15000U;
    while ((LONG)(deadline - GetTickCount()) > 0) {
        MSG message;
        int received = 0;
        while (PeekMessageA(&message, NULL, 0, 0, PM_REMOVE)) {
            received = 1;
            if (message.message == MM_WOM_OPEN) {
                ++open_messages;
            } else if (message.message == MM_WOM_DONE) {
                ++done_messages;
            } else if (message.hwnd == NULL &&
                       message.message == caller_message) {
                ++caller_messages;
            } else if (message.message == WM_QUIT) {
                quit_seen = 1;
            }
            printf("MESSAGE hwnd=%p message=0x%04lx wparam=0x%08lx "
                   "lparam=0x%08lx open=%d done=%d caller=%d\n",
                   (void *)message.hwnd, (unsigned long)message.message,
                   (unsigned long)message.wParam,
                   (unsigned long)message.lParam, open_messages, done_messages,
                   caller_messages);
            TranslateMessage(&message);
            DispatchMessageA(&message);
            query_state(get_info, "after_dispatch");
            if (quit_seen) {
                break;
            }
        }
        if (quit_seen) {
            break;
        }
        if (!received) {
            Sleep(100);
        }
        int state = query_state(get_info, "poll");
        if (state == 0 && done_messages > 0) {
            break;
        }
    }

    printf("PLAYBACK_END elapsed_ms=%lu open=%d done=%d caller=%d quit=%d "
           "state=%d\n",
           (unsigned long)(GetTickCount() - start_tick), open_messages,
           done_messages, caller_messages, quit_seen,
           query_state(get_info, "before_stop"));

    if (played) {
        stop_tts();
        query_state(get_info, "after_stop");
    }
    if (loaded) {
        unload_tts(1);
    }
    FreeLibrary(dll);
    printf("PROBE_END\n");
    return 0;
}
