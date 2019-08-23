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
    octet IV= {0,sizeof(iv),iv};
    char iv2[PQNIST_AES_IV_LENGTH];
    octet IV2= {0,sizeof(iv2),iv2};

    // Message to be sent to Bob
    char p[256];
    octet P = {0, sizeof(p), p};
    OCT_jstring(&P,"Hello Bob! This is a message from Alice");

    // Pad message
    int l = 16 - (P.len % 16);
    if (l < 16)
    {
        OCT_jbyte(&P,0,l);
    }

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
        fprintf(stderr, "FAILURE pqnist_keys rc: %d\n", rc);
        printf("FAILURE\n");
        exit(EXIT_FAILURE);
    }

    // BLS signature
    char S[SIGLEN];

    // SIKE encapsulated key
    uint8_t ek[OQS_KEM_sike_p751_length_ciphertext];

    // Alice

    printf("Alice Key: ");
    amcl_print_hex(K.val, K.len);

    // Random initialization value
    IV.len=PQNIST_AES_IV_LENGTH;
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
        fprintf(stderr, "FAILURE pqnist_encapsulate_encrypt rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    printf("Alice encrypted key: ");
    OCT_output(&K);

    // Bob

    // Obtain encapsulated AES key and decrypt K
    rc = pqnist_decapsulate_decrypt(K.val, K.len, IV2.val, SIKEsk, ek);
    if(rc)
    {
        fprintf(stderr, "FAILURE pqnist_decapsulate_decrypt rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    printf("Bob Key: ");
    amcl_print_hex(K.val, K.len);

    printf("Bob IV ");
    OCT_output(&IV);

    printf("Bob Ciphertext: ");
    OCT_output(&P);

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
        fprintf(stderr, "FAILURE OCT_comp rc: %d\n", rc);
        printf("FAILURE\n");
        exit(EXIT_FAILURE);
    }

    // Sign message

    // Alice signs message
    rc = pqnist_sign(P.val, BLSsk, S);
    if(rc)
    {
        fprintf(stderr, "FAILURE pqnist_sign rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    // Bob verifies message
    rc = pqnist_verify(P.val, BLSpk, S);
    if (rc)
    {
        fprintf(stderr, "FAILURE: verify failed!\n errorCode %d", rc);
        printf("FAILURE\n");
        exit(EXIT_FAILURE);
    }

    // clear memory
    OQS_MEM_cleanse(SIKEsk, OQS_KEM_sike_p751_length_secret_key);
    OCT_clear(&K);
    OCT_clear(&IV);
    OCT_clear(&P);
    KILL_CSPRNG(&RNG);

    printf("SUCCESS\n");
    exit(EXIT_SUCCESS);
}
