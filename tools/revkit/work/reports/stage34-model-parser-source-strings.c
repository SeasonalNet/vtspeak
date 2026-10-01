/* Ghidra pseudocode excerpt from FUN_1000d190 @ 0x1000d190. */

/* Parser rows advance by 0x94 bytes. The primary source string is read at row
 * +0x34, while the auxiliary string is read at row +0x52. The primary string
 * is copied into a 32-byte local buffer before FUN_1000d640; the auxiliary
 * string is passed as that helper's param_5. */
row = parser_rows + row_index * 0x94;
copy_c_string(local_3c, row + 0x34);
copy_c_string(local_a0, row + 0x52);
FUN_1000d640(local_3c, local_5c, final_row, parser_row[0x24], local_a0,
             row_is_special, model_context_table);

/* FUN_1000d640 pseudocode names, local sizes, and pointer relationships may
 * be imperfect. Stage 20 runtime captures directly observed +0x34 and +0x52.
 * Numeric rows contain expanded lexemes such as "one", "point", "two",
 * "five", "dollars", "percent", and "PM" while multiple rows retain one
 * original source span. The Go ordinary-row builder uses its bounded numeric
 * normalizers to populate +0x34; other number grammars remain inference.
 * Production of +0x52 remains incomplete. */
