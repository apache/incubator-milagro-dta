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
   Encapsulate a secret and use the secret to encrypt a message
   Decapsulate the secret and use the secret to decrypt the encrypted message
*/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <amcl/utils.h>
#include <amcl/randapi.h>
#include <amcl/bls_BLS381.h>
#include <oqs/oqs.h>
#include <pqnist/pqnist.h>

#define NTHREADS 8
#define MAXSIZE 256
#define G2LEN 4*BFS_BLS381

int main()
{
    int i,rc;

    // Seed value for CSPRNG
    char seed[PQNIST_SEED_LENGTH];
    octet SEED = {sizeof(seed),sizeof(seed),seed};

    // Seed value for key generation
    char seedkeys[NTHREADS][PQNIST_SEED_LENGTH];

    csprng RNG;

    // Initialization vector
    char iv[PQNIST_AES_IV_LENGTH];
    octet IV= {sizeof(iv),sizeof(iv),iv};

    // Message to be sent to Bob
    char p[NTHREADS][MAXSIZE];
    octet P[NTHREADS];

    // AES CBC ciphertext
    char c[NTHREADS][MAXSIZE];
    octet C[NTHREADS];

    // non random seed value
    for (i=0; i<32; i++) SEED.val[i]=i+1;
    printf("SEED: ");
    OCT_output(&SEED);
    printf("\n");

    // initialise random number generator
    CREATE_CSPRNG(&RNG,&SEED);

    // Initialise key generation seed
    for(i=0; i<NTHREADS; i++)
    {
        for(int j=0; j<PQNIST_SEED_LENGTH; j++)
        {
            seedkeys[i][j] = i;
        }
    }

    // Bob's SIKE keys
    uint8_t SIKEpk[NTHREADS][OQS_KEM_sike_p751_length_public_key];
    uint8_t SIKEsk[NTHREADS][OQS_KEM_sike_p751_length_secret_key];

    // Alice's BLS keys (not used)
    char BLSpk[NTHREADS][G2LEN];
    char BLSsk[NTHREADS][BGS_BLS381];

    #pragma omp parallel for
    for(i=0; i<NTHREADS; i++)
    {
        rc = pqnist_keys(seedkeys[i], SIKEpk[i], SIKEsk[i], BLSpk[i], BLSsk[i]);
        if (rc)
        {
            fprintf(stderr, "FAILURE pqnist_keys rc: %d\n", rc);
            OQS_MEM_cleanse(SIKEsk[i], OQS_KEM_sike_p751_length_secret_key);
            exit(EXIT_FAILURE);
        }

        int j = OQS_KEM_sike_p751_length_public_key;
        printf("Bob SIKE pklen %d pk: ", j);
        amcl_print_hex(SIKEpk[i], j);
        j = OQS_KEM_sike_p751_length_secret_key;
        printf("Bob SIKE sklen %d sk: ", j);
        amcl_print_hex(SIKEsk[i], j);

    }

    // Alice

    for(i=0; i<NTHREADS; i++)
    {
        bzero(p[i],sizeof(p[i]));
        P[i].max = MAXSIZE;
        P[i].len = sprintf(p[i], "Hello Bob! This is a message from Alice %d", i);
        P[i].val = p[i];
        // Pad message
        int l = 16 - (P[i].len % 16);
        if (l < 16)
        {
            OCT_jbyte(&P[i],0,l);
        }
    }

    // Random initialization value
    generateRandom(&RNG,&IV);
    printf("Alice IV: ");
    OCT_output(&IV);

    // Copy plaintext
    for(i=0; i<NTHREADS; i++)
    {
        C[i].val = c[i];
        C[i].max = MAXSIZE;
        OCT_copy(&C[i],&P[i]);
        printf("Alice Plaintext: ");
        OCT_output_string(&C[i]);
        printf("\n");
    }

    // SIKE encapsulated key
    uint8_t ek[NTHREADS][OQS_KEM_sike_p751_length_ciphertext];

    #pragma omp parallel for
    for(i=0; i<NTHREADS; i++)
    {

        // Generate an AES which is ecapsulated using SIKE. Use this key to
        // AES encrypt the K parameter.
        rc = pqnist_encapsulate_encrypt(C[i].val, C[i].len, IV.val, SIKEpk[i], ek[i]);
        if(rc)
        {
            fprintf(stderr, "FAILURE pqnist_encapsulate_encrypt rc: %d\n", rc);
            exit(EXIT_FAILURE);
        }

        printf("Alice ciphertext: ");
        OCT_output(&C[i]);

        printf("Alice ek %lu ek: ", sizeof(ek[i]));
        amcl_print_hex(ek[i], sizeof(ek[i]));
        printf("\n");

    }

    // Bob

    #pragma omp parallel for
    for(i=0; i<NTHREADS; i++)
    {
        // Obtain encapsulated AES key and decrypt C
        rc = pqnist_decapsulate_decrypt(C[i].val, C[i].len, IV.val, SIKEsk[i], ek[i]);
        if(rc)
        {
            fprintf(stderr, "FAILURE pqnist_decapsulate_decrypt rc: %d\n", rc);
            exit(EXIT_FAILURE);
        }

        printf("Bob Plaintext: ");
        OCT_output(&C[i]);

        printf("Bob Plaintext: ");
        OCT_output_string(&C[i]);
        printf("\n");

        // Compare sent and recieved message (returns 0 for failure)
        rc = OCT_comp(&P[i],&C[i]);
        if(!rc)
        {
            fprintf(stderr, "FAILURE OCT_comp rc: %d\n", rc);
            exit(EXIT_FAILURE);
        }
    }

    // clear memory

    OCT_clear(&IV);
    for(i=0; i<NTHREADS; i++)
    {
        OQS_MEM_cleanse(SIKEsk[i], OQS_KEM_sike_p751_length_secret_key);
        OCT_clear(&P[i]);
        OCT_clear(&C[i]);
    }

    KILL_CSPRNG(&RNG);

    printf("SUCCESS\n");
    exit(EXIT_SUCCESS);
}
