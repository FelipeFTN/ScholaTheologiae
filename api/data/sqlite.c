/*
 * Portuguese Accent-Insensitive Search for SQLite
 * 
 * This file allows searching for Portuguese text without caring about accents.
 * For example: searching "matrimonio" will find "matrimônio" in the database.
 * 
 * How it works:
 * 1. We have a table that maps accented characters to their basic forms
 * 2. We normalize both the search term and database text using this table  
 * 3. We then do a simple substring search on the normalized strings
 */

#include "sqlite.h"

/*
 * Portuguese Character Mapping Table
 * 
 * This table maps each character (0-255) to its "normalized" form.
 * Accented characters like À, Á, Â, Ã get mapped to 'A'
 * Regular characters like 'A', 'B', 'C' stay the same
 * Lowercase letters get converted to uppercase for case-insensitive matching
 */
static const char table_ptbr[] = {
    /* u+0000 .. U+007F - Standard ASCII */
    0x00,0x01,0x02,0x03,0x04,0x05,0x06,0x07,  0x08,0x09,0x0A,0x0B,0x0C,0x0D,0x0E,0x0F,
    0x10,0x11,0x12,0x13,0x14,0x15,0x16,0x17,  0x18,0x19,0x1A,0x1B,0x1C,0x1D,0x1E,0x1F,
    0x20,0x21,0x22,0x23,0x24,0x25,0x26,0x27,  0x28,0x29,0x2A,0x2B,0x2C,0x2D,0x2E,0x2F,
    0x30,0x31,0x32,0x33,0x34,0x35,0x36,0x37,  0x38,0x39,0x3A,0x3B,0x3C,0x3D,0x3E,0x3F,
    0x40,'A' ,'B' ,'C' ,'D' ,'E' ,'F' ,'G' ,  'H' ,'I' ,'J' ,'K' ,'L' ,'M' ,'N' ,'O' ,
    'P' ,'Q' ,'R' ,'S' ,'T' ,'U' ,'V' ,'W' ,  'X' ,'Y' ,'Z' ,0x5B,0x5C,0x5D,0x5E,0x5F,
    0x60,'A' ,'B' ,'C' ,'D' ,'E' ,'F' ,'G' ,  'H' ,'I' ,'J' ,'K' ,'L' ,'M' ,'N' ,'O' ,
    'P' ,'Q' ,'R' ,'S' ,'T' ,'U' ,'V' ,'W' ,  'X' ,'Y' ,'Z' ,0x7B,0x7C,0x7D,0x7E,0x7F,
    /* u+0080 .. U+00FF - Extended ASCII with Portuguese accent mapping */
    0x80,0x81,0x82,0x83,0x84,0x85,0x86,0x87,  0x88,0x89,0x8A,0x8B,0x8C,0x8D,0x8E,0x8F,
    0x90,0x91,0x92,0x93,0x94,0x95,0x96,0x97,  0x98,0x99,0x9A,0x9B,0x9C,0x9D,0x9E,0x9F,
    0xA0,0xA1,'C' ,0xA3,0xA4,0xA5,0xA6,0xA7,  0xA8,0xA9,'A' ,0xAB,0xAC,0xAD,0xAE,0xAF,
    0xB0,0xB1,'2' ,'3' ,0xB4,'U' ,0xB6,0xB7,  0xB8,'1' ,'O' ,0xBB,0xBC,0xBD,0xBE,0xBF,
    'A' ,'A' ,'A' ,'A' ,'A' ,'A' ,'A' ,'C' ,  'E' ,'E' ,'E' ,'E' ,'I' ,'I' ,'I' ,'I' ,
    'D' ,'N' ,'O' ,'O' ,'O' ,'O' ,'O' ,0xD7,  'O' ,'U' ,'U' ,'U' ,'U' ,'Y' ,0xDE,'S' ,
    'A' ,'A' ,'A' ,'A' ,'A' ,'A' ,'A' ,'C' ,  'E' ,'E' ,'E' ,'E' ,'I' ,'I' ,'I' ,'I' ,
    'D' ,'N' ,'O' ,'O' ,'O' ,'O' ,'O' ,0xF7,  'O' ,'U' ,'U' ,'U' ,'U' ,'Y' ,0xFE,'Y'
};

/*
 * STEP 1: Convert accented characters to their base forms
 * 
 * This function takes a string like "matrimônio" and converts it to "MATRIMONIO"
 * - Removes accents: ô -> O
 * - Converts to uppercase: m -> M (for case-insensitive matching)
 * 
 * Parameters:
 *   input: The original string (e.g., "matrimônio")
 *   input_len: Length of the input string
 *   output: Buffer where the normalized string will be stored
 *   output_len: Size of the output buffer (gets updated with actual length)
 */
static void normalize_string(const char* input, int input_len, char* output, int* output_len) {
    // Set up pointers to walk through the input string
    const unsigned char* current_char = (const unsigned char*)input;
    const unsigned char* end_of_input = current_char + input_len;
    char* output_pos = output;
    int max_output_space = *output_len - 1; // Leave room for null terminator
    
    // Process each character in the input string
    while (current_char < end_of_input && (output_pos - output) < max_output_space) {
        unsigned char c = *current_char++;
        
        if (c < 0x80) {
            // Regular ASCII character (a-z, A-Z, 0-9, etc.)
            // Use our mapping table to normalize it (e.g., 'a' -> 'A', 'ã' -> 'A')
            *output_pos++ = table_ptbr[c];
        } else if ((c & 0xFE) == 0xC2 && current_char < end_of_input) {
            // UTF-8 accented character (like ã, é, ô, etc.)
            // These are encoded as two bytes in UTF-8
            unsigned char second_byte = *current_char++;
            unsigned char mapped_char = table_ptbr[0x80 | ((c << 6) & 0x40) | (second_byte & 0x3F)];
            *output_pos++ = mapped_char;
        } else {
            // Other characters - just copy them as-is
            *output_pos++ = c;
        }
    }
    
    // Null-terminate the string and update the length
    *output_pos = '\0';
    *output_len = output_pos - output;
}

/*
 * STEP 2: Check if one string contains another (simple substring search)
 * 
 * This function checks if 'pattern' is found anywhere inside 'text'
 * Both strings should already be normalized (accents removed, uppercase)
 * 
 * Example: 
 *   text = "MATRIMONIO", pattern = "MATRI" -> returns 1 (found)
 *   text = "MATRIMONIO", pattern = "CASA" -> returns 0 (not found)
 * 
 * Parameters:
 *   text: The text to search in (already normalized)
 *   text_len: Length of the text
 *   pattern: The pattern to search for (already normalized)
 *   pattern_len: Length of the pattern
 * 
 * Returns: 1 if pattern is found in text, 0 if not found
 */
static int contains_pattern(const char* text, int text_len, const char* pattern, int pattern_len) {
    // Empty pattern matches everything
    if (pattern_len == 0) {
        return 1;
    }
    
    // If text is shorter than pattern, it can't contain the pattern
    if (text_len < pattern_len) {
        return 0;
    }
    
    // Try each possible starting position in the text
    for (int start_pos = 0; start_pos <= text_len - pattern_len; start_pos++) {
        // Check if pattern matches at this position
        int matches = 1; // Assume it matches until proven otherwise
        
        for (int i = 0; i < pattern_len; i++) {
            if (text[start_pos + i] != pattern[i]) {
                matches = 0; // Found a mismatch
                break;
            }
        }
        
        if (matches) {
            return 1; // Found a match!
        }
    }
    
    return 0; // Pattern not found anywhere in text
}

/*
 * MAIN FUNCTION: SQLite calls this function for accent_insensitive_like()
 * 
 * This is the function that SQLite calls when you use:
 * SELECT * FROM table WHERE accent_insensitive_like(column, 'search_term')
 * 
 * It combines STEP 1 and STEP 2 above:
 * 1. Normalizes both the database text and the search pattern
 * 2. Checks if the normalized pattern is found in the normalized text
 * 
 * Parameters:
 *   context: SQLite function context (for returning results)
 *   argc: Number of arguments (should be 2)
 *   argv: Array of arguments [database_text, search_pattern]
 */
static void accent_insensitive_like(sqlite3_context *context, int argc, sqlite3_value **argv) {
    // Make sure we got exactly 2 arguments
    if (argc != 2) {
        sqlite3_result_error(context, "accent_insensitive_like() needs exactly 2 arguments", -1);
        return;
    }
    
    // Get the text from the database column and the search pattern
    const char* database_text = (const char*)sqlite3_value_text(argv[0]);
    const char* search_pattern = (const char*)sqlite3_value_text(argv[1]);
    
    // Handle null values
    if (!database_text || !search_pattern) {
        sqlite3_result_int(context, 0); // Return "no match" for null values
        return;
    }
    
    // Get the lengths of both strings
    int text_length = sqlite3_value_bytes(argv[0]);
    int pattern_length = sqlite3_value_bytes(argv[1]);
    
    // Allocate memory for the normalized versions
    // Make the buffers a bit larger than needed for safety
    int buffer_size = text_length + pattern_length + 20;
    char* normalized_text = (char*)malloc(buffer_size);
    char* normalized_pattern = (char*)malloc(buffer_size);
    
    // Check if memory allocation failed
    if (!normalized_text || !normalized_pattern) {
        free(normalized_text);
        free(normalized_pattern);
        sqlite3_result_error(context, "Out of memory", -1);
        return;
    }
    
    // Normalize both strings (remove accents, convert to uppercase)
    int norm_text_len = buffer_size;
    int norm_pattern_len = buffer_size;
    normalize_string(database_text, text_length, normalized_text, &norm_text_len);
    normalize_string(search_pattern, pattern_length, normalized_pattern, &norm_pattern_len);
    
    // Check if the normalized pattern is found in the normalized text
    int found = contains_pattern(normalized_text, norm_text_len, normalized_pattern, norm_pattern_len);
    
    // Clean up memory
    free(normalized_text);
    free(normalized_pattern);
    
    // Return the result to SQLite (1 for match, 0 for no match)
    sqlite3_result_int(context, found);
}

/*
 * SETUP FUNCTION: Registers our function with SQLite
 * 
 * This tells SQLite about our new function so it can be used in SQL queries
 */
int setup_accent_insensitive_like(sqlite3 *db) {
    if (db == NULL) {
        return 1; // Error: no database connection
    }
    
    // Register our function with SQLite
    int result = sqlite3_create_function(
        db,                           // Database connection
        "accent_insensitive_like",    // Function name (used in SQL)
        2,                           // Number of arguments the function takes
        SQLITE_UTF8,                 // Text encoding
        NULL,                        // No user data needed
        accent_insensitive_like,     // Pointer to our function
        NULL,                        // Not an aggregate function
        NULL                         // Not an aggregate function
    );
    
    if (result != SQLITE_OK) {
        return 1; // Error registering function
    }
    
    return 0; // Success
}

/*
 * WRAPPER FUNCTIONS: These allow Go to call our C functions directly
 * 
 * Go can't call our static functions directly, so we provide these
 * wrapper functions that Go can access.
 */

// Wrapper to normalize a string (Go -> C)
void normalize_string_wrapper(const char* input, int input_len, char* output, int* output_len) {
    normalize_string(input, input_len, output, output_len);
}

// Wrapper to check if text contains pattern (Go -> C) 
int contains_pattern_wrapper(const char* text, int text_len, const char* pattern, int pattern_len) {
    return contains_pattern(text, text_len, pattern, pattern_len);
}
