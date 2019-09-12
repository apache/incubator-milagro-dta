/*
	Licensed to the Apache Software Foundation (ASF) under one
	or more contributor license agreements.  See the NOTICE file
	distributed with this work for additional information
	regarding copyright ownership.  The ASF licenses this file
	to you under the Apache License, Version 2.0 (the
	"License"); you may not use this file except in compliance
	with the License.  You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing,
	software distributed under the License is distributed on an
	"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
	KIND, either express or implied.  See the License for the
	specific language governing permissions and limitations
	under the License.
*/

/**
 * @file run_sike.c
 * @author Kealan McCusker
 * @brief Encapsulate and decapsulate a secret using SIKE
 */

#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <oqs/oqs.h>

// Set to have fixed values for the key pair
#define FIXED

/* Print encoded binary string in hex */
static void printHex(uint8_t *src, int src_len)
{
    int i;
    for (i = 0; i < src_len; i++)
    {
        printf("%02x", (unsigned char) src[i]);
    }
    printf("\n");
}

/* Cleaning up memory etc */
static void cleanup(uint8_t *secret_key, size_t secret_key_len,
                    uint8_t *shared_secret_e, uint8_t *shared_secret_d,
                    size_t shared_secret_len);

/* This function gives an example of the operations performed by both
 * the decapsulator and the encapsulator in a single KEM session,
 * using only compile-time macros and allocating variables
 * statically on the stack, calling a specific algorithm's functions
 * directly.
 */
static OQS_STATUS example()
{
    uint8_t public_key[OQS_KEM_sike_p751_length_public_key];
    uint8_t secret_key[OQS_KEM_sike_p751_length_secret_key];
    uint8_t ciphertext[OQS_KEM_sike_p751_length_ciphertext];
    uint8_t shared_secret_e[OQS_KEM_sike_p751_length_shared_secret];
    uint8_t shared_secret_d[OQS_KEM_sike_p751_length_shared_secret];

#ifdef FIXED
    uint8_t entropy_input[48];
    for (size_t i = 0; i < 48; i++)
    {
        entropy_input[i] = i;
    }

    OQS_STATUS rc = OQS_randombytes_switch_algorithm(OQS_RAND_alg_nist_kat);
    if (rc != OQS_SUCCESS)
    {
        return rc;
    }
    OQS_randombytes_nist_kat_init(entropy_input, NULL, 256);
#endif

    rc = OQS_KEM_sike_p751_keypair(public_key, secret_key);
    if (rc != OQS_SUCCESS)
    {
        fprintf(stderr, "ERROR: OQS_KEM_sike_p751_keypair failed!\n");
        cleanup(secret_key, OQS_KEM_sike_p751_length_secret_key,
                shared_secret_e, shared_secret_d,
                OQS_KEM_sike_p751_length_shared_secret);

        return OQS_ERROR;
    }
    int i = OQS_KEM_sike_p751_length_public_key;
    printf("pklen %d pk: ", i);
    printHex(public_key, i);
    i = OQS_KEM_sike_p751_length_secret_key;
    printf("sklen %d sk: ", i);
    printHex(secret_key, i);

    rc = OQS_KEM_sike_p751_encaps(ciphertext, shared_secret_e, public_key);
    if (rc != OQS_SUCCESS)
    {
        fprintf(stderr, "ERROR: OQS_KEM_sike_p751_encaps failed!\n");
        cleanup(secret_key, OQS_KEM_sike_p751_length_secret_key,
                shared_secret_e, shared_secret_d,
                OQS_KEM_sike_p751_length_shared_secret);

        return OQS_ERROR;
    }
    i = OQS_KEM_sike_p751_length_ciphertext;
    printf("ciphertextlen %d ciphertext: ", i);
    printHex(ciphertext, i);

    rc = OQS_KEM_sike_p751_decaps(shared_secret_d, ciphertext, secret_key);
    if (rc != OQS_SUCCESS)
    {
        fprintf(stderr, "ERROR: OQS_KEM_sike_p751_decaps failed!\n");
        cleanup(secret_key, OQS_KEM_sike_p751_length_secret_key,
                shared_secret_e, shared_secret_d,
                OQS_KEM_sike_p751_length_shared_secret);

        return OQS_ERROR;
    }
    i = OQS_KEM_sike_p751_length_shared_secret;
    printf("shared_secret_elen %d shared_secret_e: ", i);
    printHex(shared_secret_e, i);
    printf("shared_secret_dlen %d shared_secret_d: ", i);
    printHex(shared_secret_d, i);

    printf("OQS_KEM_sike_p751 operations completed.\n");

    return OQS_SUCCESS; // success!
}


int main(void)
{
    example();
}

static void cleanup(uint8_t *secret_key, size_t secret_key_len,
                    uint8_t *shared_secret_e, uint8_t *shared_secret_d,
                    size_t shared_secret_len)
{
    OQS_MEM_cleanse(secret_key, secret_key_len);
    OQS_MEM_cleanse(shared_secret_e, shared_secret_len);
    OQS_MEM_cleanse(shared_secret_d, shared_secret_len);
}

