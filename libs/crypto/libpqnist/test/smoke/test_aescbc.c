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
   Test AES-256 encryption and decryption in CBC mode
*/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <amcl/utils.h>
#include <amcl/randapi.h>
#include <amcl/pqnist.h>

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
    char p1[256];
    octet P1 = {0, sizeof(p1), p1};
    OCT_jstring(&P1,"Hello Bob! This is a message from Alice");

    // Recovered plaintext
    char p2[256];
    octet P2= {0,sizeof(p2),p2};

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
    OCT_output(&P1);

    printf("Alice Plaintext: ");
    OCT_output_string(&P1);
    printf("\n");

    // Encrypt plaintext
    pqnist_aes_cbc_encrypt(K.val, K.len, IV.val, P1.val, P1.len);

    printf("Alice: Ciphertext: ");
    OCT_output(&P1);

    // Bob

    printf("Bob Key: ");
    amcl_print_hex(K.val, K.len);

    printf("Bob IV ");
    OCT_output(&IV);

    OCT_copy(&P2,&P1);
    printf("Bob Ciphertext: ");
    OCT_output(&P2);

    pqnist_aes_cbc_decrypt(K.val, K.len, IV.val, P2.val, P2.len);

    printf("Bob Plaintext: ");
    OCT_output(&P2);

    printf("Bob Plaintext: ");
    OCT_output_string(&P2);
    printf("\n");

    // Expected message
    char p[256];
    octet P = {0, sizeof(p), p};
    OCT_jstring(&P,"Hello Bob! This is a message from Alice");

    if (!OCT_comp(&P,&P2))
    {
        printf("FAILURE Decryption\n");
        printf("P: ");
        OCT_output(&P);
        printf("P2: ");
        OCT_output(&P2);
        exit(EXIT_FAILURE);
    }

    /* clear memory */
    OCT_clear(&K);
    OCT_clear(&IV);
    OCT_clear(&P1);
    OCT_clear(&P2);

    KILL_CSPRNG(&RNG);

    printf("SUCCESS\n");
    exit(EXIT_SUCCESS);
}
