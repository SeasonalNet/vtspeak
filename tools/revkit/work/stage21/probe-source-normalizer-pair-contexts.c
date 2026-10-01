#include <windows.h>

#include <stdint.h>
#include <stdio.h>
#include <string.h>

typedef short(__cdecl *source_normalizer_fn)(char *output, char *input);

static int is_trim_byte(unsigned char value) {
  return value == 0x20 || value == 0x09 || value == 0x0a || value == 0x0d;
}

static int is_rejected_pair(unsigned char first, unsigned char second) {
  return ((first >= 0xa1 && first <= 0xad && second >= 0xa1 && second <= 0xfe) ||
          (first == 0xae && second >= 0xa1 && second <= 0xc2) ||
          (first == 0xfd && second == 0xfe));
}

static uint64_t hash_byte(uint64_t hash, unsigned char value) {
  return (hash ^ value) * UINT64_C(1099511628211);
}

static int call_and_check(source_normalizer_fn normalize,
                          unsigned char input[4], unsigned int *rejected,
                          unsigned int *empty, unsigned int *one,
                          unsigned int *two, unsigned int *three,
                          unsigned int *mismatches, uint64_t *hash) {
  unsigned char original[4];
  unsigned char output[52];
  memcpy(original, input, sizeof(original));
  memset(output, 0xcc, sizeof(output));
  output[0] = 0;

  short result = normalize((char *)output, (char *)input);
  unsigned int begin = 0;
  unsigned int end = 0;
  while (end < 3 && input[end] != 0) {
    ++end;
  }
  while (begin < end && is_trim_byte(input[begin])) {
    ++begin;
  }
  while (end > begin && is_trim_byte(input[end - 1])) {
    --end;
  }
  unsigned int expected_length = end - begin;
  int expected_rejection = 0;
  unsigned int rejection_index = end;
  for (unsigned int index = begin; index + 1 < end; ++index) {
    if (is_rejected_pair(input[index], input[index + 1])) {
      expected_rejection = 1;
      rejection_index = index;
      break;
    }
  }

  if (expected_rejection) {
    ++*rejected;
    unsigned int partial_length = rejection_index - begin;
    if (result != -3 ||
        memcmp(output, input + begin, partial_length) != 0 ||
        (partial_length == 0 && output[0] != 0) ||
        (partial_length > 0 && output[partial_length] != 0xcc)) {
      if (*mismatches < 8) {
        printf("SOURCE_PAIR_CONTEXT_MISMATCH input=%02x,%02x,%02x result=%d "
               "output0=%02x expected_rejection=1\n",
               original[0], original[1], original[2], result, output[0]);
      }
      ++*mismatches;
    }
  } else if (expected_length == 0) {
    ++*empty;
    if (result != -1 || output[0] != 0) {
      ++*mismatches;
    }
  } else {
    if (result != (int)expected_length ||
        memcmp(output, input + begin, expected_length) != 0 ||
        output[expected_length] != 0) {
      ++*mismatches;
    }
    if (expected_length == 1) {
      ++*one;
    } else if (expected_length == 2) {
      ++*two;
    } else {
      ++*three;
    }
  }
  if (memcmp(input, original, sizeof(original)) != 0) {
    ++*mismatches;
  }

  for (unsigned int index = 0; index < 3; ++index) {
    *hash = hash_byte(*hash, original[index]);
  }
  uint32_t result_bits = (uint32_t)result;
  for (unsigned int index = 0; index < 4; ++index) {
    *hash = hash_byte(*hash,
                      (unsigned char)(result_bits >> (index * 8)));
  }
  unsigned int hash_length = result > 0 ? (unsigned int)result : 1;
  for (unsigned int index = 0; index < hash_length; ++index) {
    *hash = hash_byte(*hash, output[index]);
  }
  return 0;
}

int main(void) {
  HMODULE module = LoadLibraryA("Z:\\samples\\vt_pau.dll");
  if (module == NULL) {
    printf("SOURCE_PAIR_CONTEXT_ERROR LoadLibraryA=%lu\n", GetLastError());
    return 1;
  }
  source_normalizer_fn normalize =
      (source_normalizer_fn)((unsigned char *)module + 0x5f2e0);

  unsigned int total_calls = 0;
  unsigned int total_rejected = 0;
  unsigned int total_empty = 0;
  unsigned int total_one = 0;
  unsigned int total_two = 0;
  unsigned int total_three = 0;
  unsigned int total_mismatches = 0;
  uint64_t total_hash = UINT64_C(14695981039346656037);

  for (unsigned int placement = 0; placement < 2; ++placement) {
    unsigned int calls = 0;
    unsigned int rejected = 0;
    unsigned int empty = 0;
    unsigned int one = 0;
    unsigned int two = 0;
    unsigned int three = 0;
    unsigned int mismatches = 0;
    uint64_t hash = UINT64_C(14695981039346656037);

    for (unsigned int first_pair_byte = 0; first_pair_byte <= 0xff;
         ++first_pair_byte) {
      for (unsigned int second_pair_byte = 0; second_pair_byte <= 0xff;
           ++second_pair_byte) {
        if (!is_rejected_pair((unsigned char)first_pair_byte,
                              (unsigned char)second_pair_byte)) {
          continue;
        }
        for (unsigned int context_byte = 0; context_byte <= 0xff;
             ++context_byte) {
          unsigned char input[4] = {0, 0, 0, 0};
          if (placement == 0) {
            input[0] = (unsigned char)context_byte;
            input[1] = (unsigned char)first_pair_byte;
            input[2] = (unsigned char)second_pair_byte;
          } else {
            input[0] = (unsigned char)first_pair_byte;
            input[1] = (unsigned char)second_pair_byte;
            input[2] = (unsigned char)context_byte;
          }
          call_and_check(normalize, input, &rejected, &empty, &one, &two,
                         &three, &mismatches, &hash);
          ++calls;
        }
      }
    }

    total_calls += calls;
    total_rejected += rejected;
    total_empty += empty;
    total_one += one;
    total_two += two;
    total_three += three;
    total_mismatches += mismatches;
    for (unsigned int index = 0; index < 8; ++index) {
      total_hash = hash_byte(total_hash, (unsigned char)(hash >> (index * 8)));
    }
    printf("SOURCE_PAIR_CONTEXT placement=%s calls=%u rejected=%u empty=%u "
           "one=%u two=%u three=%u mismatches=%u fnv64=%016llx\n",
           placement == 0 ? "prefix" : "suffix", calls, rejected, empty,
           one, two, three, mismatches, (unsigned long long)hash);
  }

  printf("SOURCE_PAIR_CONTEXT_MATRIX_COMPLETE calls=%u rejected=%u empty=%u "
         "one=%u two=%u three=%u mismatches=%u fnv64=%016llx\n",
         total_calls, total_rejected, total_empty, total_one, total_two,
         total_three, total_mismatches, (unsigned long long)total_hash);
  FreeLibrary(module);
  return total_calls == 643584 && total_mismatches == 0 ? 0 : 2;
}
