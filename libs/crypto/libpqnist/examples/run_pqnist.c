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
   Run through the flow of encrypting, ecapsulating and signing a message
*/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <amcl/utils.h>
#include <amcl/randapi.h>
#include <amcl/bls_BLS381.h>
#include <oqs/oqs.h>
#include <pqnist/pqnist.h>

#define G2LEN 4*BFS_BLS381
#define SIGLEN BFS_BLS381+1

int main()
{
    int i,rc;

    // Seed value for CSPRNG
    char seed[PQNIST_SEED_LENGTH];
    octet SEED = {sizeof(seed),sizeof(seed),seed};

    csprng RNG;

    // AES Key
    char k[PQNIST_AES_KEY_LENGTH];
    octet K= {0,sizeof(k),k};

    // Initialization vectors
    char iv[PQNIST_AES_IV_LENGTH];
    octet IV= {sizeof(iv),sizeof(iv),iv};
    char iv2[PQNIST_AES_IV_LENGTH];
    octet IV2= {sizeof(iv2),sizeof(iv2),iv2};

    // Message to be sent to Bob
    char p[256];
    octet P = {0, sizeof(p), p};
    OCT_jstring(&P,"Hello Bob! This is a message from Alice");

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

    // AES CBC ciphertext
    char c[256];
    octet C = {0, sizeof(c), c};

    // non random seed value
    for (i=0; i<PQNIST_SEED_LENGTH; i++) SEED.val[i]=i+1;
    printf("SEED: ");
    OCT_output(&SEED);
    printf("\n");

    // initialise random number generator
    CREATE_CSPRNG(&RNG,&SEED);

    // Generate 256 bit AES Key
    K.len=PQNIST_AES_KEY_LENGTH;
    generateRandom(&RNG,&K);

    // Generate SIKE and BLS keys

    // Bob's SIKE keys
    uint8_t SIKEpk[OQS_KEM_sike_p751_length_public_key];
    uint8_t SIKEsk[OQS_KEM_sike_p751_length_secret_key];

    // Alice's BLS keys
    char BLSsk[BGS_BLS381];
    char BLSpk[G2LEN];

    rc = pqnist_keys(seed, SIKEpk, SIKEsk, BLSpk, BLSsk);
    if (rc)
    {
        fprintf(stderr, "ERROR pqnist_keys rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    i = OQS_KEM_sike_p751_length_public_key;
    printf("Bob SIKE pklen %d pk: ", i);
    amcl_print_hex(SIKEpk, i);
    i = OQS_KEM_sike_p751_length_secret_key;
    printf("Bob SIKE sklen %d sk: ", i);
    amcl_print_hex(SIKEsk, i);
    printf("BLS pklen %d pk: ", G2LEN);
    amcl_print_hex(BLSpk, G2LEN);
    printf("BLS sklen %d BLS sk: ", BGS_BLS381);
    amcl_print_hex(BLSsk, BGS_BLS381);
    printf("\n");

    // BLS signature
    char S[SIGLEN];

    // SIKE encapsulated key
    uint8_t ek[OQS_KEM_sike_p751_length_ciphertext];

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

    // Copy plaintext
    OCT_copy(&C,&P);

    printf("Alice Plaintext: ");
    OCT_output_string(&C);
    printf("\n");

    // Encrypt plaintext
    pqnist_aes_cbc_encrypt(K.val, K.len, IV.val, C.val, C.len);

    printf("Alice Ciphertext: ");
    OCT_output(&C);

    generateRandom(&RNG,&IV2);
    printf("Alice IV2: ");
    OCT_output(&IV2);

    // Generate an AES which is ecapsulated using SIKE. Use this key to
    // AES encrypt the K parameter.
    rc = pqnist_encapsulate_encrypt(K.val, K.len, IV2.val, SIKEpk, ek);
    if(rc)
    {
        fprintf(stderr, "ERROR pqnist_encapsulate_encrypt rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    printf("Alice encrypted key: ");
    OCT_output(&K);

    i = OQS_KEM_sike_p751_length_ciphertext;
    printf("Alice ek1 %d ek1: ", i);
    amcl_print_hex(ek, i);
    printf("Alice ek2 %d ek2: ", i);
    amcl_print_hex(&ek[OQS_KEM_sike_p751_length_ciphertext], i);
    printf("\n");

    // Bob

    // Obtain encapsulated AES key and decrypt K
    rc = pqnist_decapsulate_decrypt(K.val, K.len, IV2.val, SIKEsk, ek);
    if(rc)
    {
        fprintf(stderr, "ERROR pqnist_decapsulate_decrypt rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    printf("Bob Key: ");
    amcl_print_hex(K.val, K.len);

    printf("Bob IV ");
    OCT_output(&IV);

    printf("Bob Ciphertext: ");
    OCT_output(&C);

    pqnist_aes_cbc_decrypt(K.val, K.len, IV.val, C.val, C.len);

    printf("Bob Plaintext: ");
    OCT_output(&C);

    printf("Bob Plaintext: ");
    OCT_output_string(&C);
    printf("\n");

    // Compare sent and recieved message (returns 0 for failure)
    rc = OCT_comp(&P,&C);
    if(!rc)
    {
        fprintf(stderr, "ERROR OCT_comp rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }


    // Sign message

    // Alice signs message
    rc = pqnist_sign(P.val, BLSsk, S);
    if(rc)
    {
        fprintf(stderr, "ERROR pqnist_sign rc: %d\n", rc);
        printf("FAILURE\n");
        exit(EXIT_FAILURE);
    }

    printf("Alice Slen %d SIG", SIGLEN);
    amcl_print_hex(S, SIGLEN);
    printf("\n");

    // Bob verifies message
    rc = pqnist_verify(P.val, BLSpk, S);
    if (rc == BLS_OK)
    {
        printf("BOB SUCCESS: signature verified\n");
    }
    else
    {
        fprintf(stderr, "BOB ERROR: verify failed!\n errorCode %d", rc);
        exit(EXIT_FAILURE);
    }


    printf("Bob P ");
    OCT_output(&P);
    printf("\n");

    // Bob verifies corrupted message
    char tmp = P.val[0];
    P.val[0] = 0;
    rc = pqnist_verify(P.val, BLSpk, S);
    if (rc == BLS_OK)
    {
        printf("BOB SUCCESS: signature verified\n");
    }
    else
    {
        fprintf(stderr, "BOB ERROR verify failed! errorCode: %d\n", rc);
    }

    // Fix message
    P.val[0] = tmp;
    printf("Bob P ");
    OCT_output(&P);
    printf("\n");

    // Check signature is correct
    rc = pqnist_verify(P.val, BLSpk, S);
    if (rc == BLS_OK)
    {
        printf("BOB SUCCESS: signature verified\n");
    }
    else
    {
        fprintf(stderr, "BOB ERROR verify failed! errorCode: %d\n", rc);
    }

    // Bob verifies corrupted signature
    S[0] = 0;
    rc = pqnist_verify(P.val, BLSpk, S);
    if (rc == BLS_OK)
    {
        printf("BOB SUCCESS: signature verified\n");
    }
    else
    {
        fprintf(stderr, "BOB ERROR verify failed! errorCode: %d\n", rc);
    }

    // clear memory
    OQS_MEM_cleanse(SIKEsk, OQS_KEM_sike_p751_length_secret_key);
    OQS_MEM_cleanse(BLSsk, OQS_SIG_picnic_L5_FS_length_secret_key);
    OCT_clear(&K);
    OCT_clear(&IV);
    OCT_clear(&P);

    KILL_CSPRNG(&RNG);

    exit(EXIT_SUCCESS);
}
