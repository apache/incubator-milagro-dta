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

/*
   Run AES-256 encryption and decryption in CBC mode
*/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <amcl/utils.h>
#include <amcl/randapi.h>
#include <pqnist/pqnist.h>

int main()
{
    int i;

    // Seed value for CSPRNG
    char seed[PQNIST_SEED_LENGTH];
    octet SEED = {sizeof(seed),sizeof(seed),seed};

    csprng RNG;

    // AES Key
    char k[PQNIST_AES_KEY_LENGTH];
    octet K= {sizeof(k),sizeof(k),k};

    // Initialization vector
    char iv[PQNIST_AES_IV_LENGTH];
    octet IV= {sizeof(iv),sizeof(iv),iv};

    // Message to be sent to Bob
    char p[256];
    octet P = {0, sizeof(p), p};
    // OCT_jstring(&P,"Hello Bob!");
    OCT_jstring(&P,"Hello Bob! This is a message from Alice");

    // non random seed value
    for (i=0; i<PQNIST_SEED_LENGTH; i++) SEED.val[i]=i+1;
    printf("SEED: ");
    OCT_output(&SEED);
    printf("\n");

    // initialise random number generator
    CREATE_CSPRNG(&RNG,&SEED);

    // Generate 256 bit AES Key
    generateRandom(&RNG,&K);

    // Alice

    printf("Alice Key: ");
    amcl_print_hex(K.val, K.len);

    // Random initialization value
    generateRandom(&RNG,&IV);
    printf("Alice IV: ");
    OCT_output(&IV);

    printf("Alice Plaintext: ");
    OCT_output(&P);

    printf("Alice Plaintext: ");
    OCT_output_string(&P);
    printf("\n");
    printf("Alice Pliantext hex:");
    OCT_output(&P);

    printf("PLAINTEXTLen = %d blocks %0.2f \n", P.len, (float) P.len/16);

    // Pad message
    int l = 16 - (P.len % 16);
    if (l < 16)
    {
        OCT_jbyte(&P,0,l);
    }

    printf("Alice Plaintext: ");
    OCT_output_string(&P);
    printf("\n");
    printf("Alice Pliantext hex:");
    OCT_output(&P);

    printf("PLAINTEXTLen = %d blocks %0.2f \n", P.len, (float) P.len/16);

    // Encrypt plaintext
    pqnist_aes_cbc_encrypt(K.val, K.len, IV.val, P.val, P.len);

    printf("Alice: Ciphertext: ");
    OCT_output(&P);

    // Bob

    printf("Bob Key: ");
    amcl_print_hex(K.val, K.len);

    printf("Bob IV ");
    OCT_output(&IV);

    printf("Bob Ciphertext: ");
    OCT_output(&P);

    pqnist_aes_cbc_decrypt(K.val, K.len, IV.val, P.val, P.len);

    printf("Bob Plaintext: ");
    OCT_output(&P);

    printf("Bob Plaintext: ");
    OCT_output_string(&P);
    printf("\n");

    /* clear memory */
    OCT_clear(&K);
    OCT_clear(&IV);
    OCT_clear(&P);

    KILL_CSPRNG(&RNG);

    return 0;
}
