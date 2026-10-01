#include <windows.h>

#include <stdint.h>
#include <stdio.h>
#include <string.h>

typedef int(__cdecl *source_normalizer_fn)(char *output, char *input);

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

int main(void) {
  HMODULE module = LoadLibraryA("Z:\\samples\\vt_pau.dll");
  if (module == NULL) {
    printf("SOURCE_PAIR_MATRIX_ERROR LoadLibraryA=%lu\n", GetLastError());
    return 1;
  }

  source_normalizer_fn normalize =
      (source_normalizer_fn)((unsigned char *)module + 0x5f2e0);
  unsigned int total_calls = 0;
  unsigned int total_rejected = 0;
  unsigned int total_empty = 0;
  unsigned int total_one = 0;
  unsigned int total_two = 0;
  unsigned int total_mismatch = 0;
  uint64_t total_hash = UINT64_C(14695981039346656037);

  for (unsigned int first_value = 0; first_value <= 0xff; ++first_value) {
    unsigned int rejected = 0;
    unsigned int empty = 0;
    unsigned int one = 0;
    unsigned int two = 0;
    unsigned int mismatch = 0;
    uint64_t row_hash = UINT64_C(14695981039346656037);

    for (unsigned int second_value = 0; second_value <= 0xff; ++second_value) {
      unsigned char input[3] = {(unsigned char)first_value,
                                (unsigned char)second_value, 0};
      unsigned char output[52];
      memset(output, 0xcc, sizeof(output));
      output[0] = 0;

      int result = normalize((char *)output, (char *)input);
      unsigned int begin = 0;
      unsigned int end = first_value == 0 ? 0 : (second_value == 0 ? 1 : 2);
      while (begin < end && is_trim_byte(input[begin])) {
        ++begin;
      }
      while (end > begin && is_trim_byte(input[end - 1])) {
        --end;
      }
      unsigned int expected_length = end - begin;
      int expected_rejection = expected_length == 2 &&
                               is_rejected_pair(input[begin], input[begin + 1]);

      if (expected_rejection) {
        if (result != -3 || output[0] != 0) {
          ++mismatch;
        }
      } else if (expected_length == 0) {
        if (result != -1 || output[0] != 0) {
          ++mismatch;
        }
        ++empty;
      } else {
        if (result != (int)expected_length ||
            memcmp(output, input + begin, expected_length) != 0 ||
            output[expected_length] != 0) {
          ++mismatch;
        }
        if (expected_length == 1) {
          ++one;
        } else {
          ++two;
        }
      }
      if (input[0] != (unsigned char)first_value ||
          input[1] != (unsigned char)second_value || input[2] != 0) {
        ++mismatch;
      }

      if (expected_rejection) {
        ++rejected;
      }
      row_hash = hash_byte(row_hash, (unsigned char)first_value);
      row_hash = hash_byte(row_hash, (unsigned char)second_value);
      uint32_t result_bits = (uint32_t)result;
      for (unsigned int byte_index = 0; byte_index < 4; ++byte_index) {
        row_hash = hash_byte(row_hash,
                             (unsigned char)(result_bits >> (byte_index * 8)));
      }
      unsigned int hash_length = result > 0 ? (unsigned int)result : 1;
      for (unsigned int byte_index = 0; byte_index < hash_length; ++byte_index) {
        row_hash = hash_byte(row_hash, output[byte_index]);
      }
      ++total_calls;
    }

    total_rejected += rejected;
    total_empty += empty;
    total_one += one;
    total_two += two;
    total_mismatch += mismatch;
    for (unsigned int byte_index = 0; byte_index < 8; ++byte_index) {
      total_hash = hash_byte(total_hash,
                             (unsigned char)(row_hash >> (byte_index * 8)));
    }
    printf("SOURCE_PAIR_ROW first=0x%02x calls=256 rejected=%u empty=%u "
           "one=%u two=%u mismatches=%u fnv64=%016llx\n",
           first_value, rejected, empty, one, two, mismatch,
           (unsigned long long)row_hash);
  }

  printf("SOURCE_PAIR_MATRIX_COMPLETE calls=%u rejected=%u empty=%u one=%u "
         "two=%u mismatches=%u fnv64=%016llx\n",
         total_calls, total_rejected, total_empty, total_one, total_two,
         total_mismatch, (unsigned long long)total_hash);
  FreeLibrary(module);
  return total_calls == 65536 && total_mismatch == 0 ? 0 : 2;
}
