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

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <oqs/oqs.h>
#include <amcl/pqnist.h>
#include <amcl/utils.h>
#include <amcl/randapi.h>
#include <amcl/bls_BLS381.h>

#define G2LEN 4*BFS_BLS381
#define SIGLEN BFS_BLS381+1

//  Generate BLS public and private key
int pqnist_bls_keys(char* seed, char* BLSpk, char* BLSsk)
{
    int rc;

    octet PK = {0,G2LEN,BLSpk};
    octet SK = {0,BGS_BLS381,BLSsk};

    if (seed==NULL)
    {
        rc =  BLS_BLS381_KEY_PAIR_GENERATE(NULL,&SK,&PK);
    }
    else
    {
        octet SEED = {PQNIST_SEED_LENGTH,PQNIST_SEED_LENGTH,seed};

        // CSPRNG
        csprng RNG;

        // initialise strong RNG
        CREATE_CSPRNG(&RNG,&SEED);

        // Generate BLS key pair
        rc =  BLS_BLS381_KEY_PAIR_GENERATE(&RNG,&SK,&PK);
    }

    if (rc)
    {
        return rc;
    }

#ifdef DEBUG
    printf("pqnist_keys BLSpklen %d BLSpk ", PK.len);
    OCT_output(&PK);
    printf("pqnist_keys BLSsklen %d BLSsk ", SK.len);
    OCT_output(&SK);
#endif

    return 0;
}


//  Sign message using the BLS algorithm
int pqnist_bls_sign(char* M, char* sk, char* S)
{
    octet SIG = {0,SIGLEN,S};
    octet SK = {BGS_BLS381,BGS_BLS381,sk};

#ifdef DEBUG
    printf("pqnist_bls_sign SK: ");
    OCT_output(&SK);
#endif

    int rc = BLS_BLS381_SIGN(&SIG,M,&SK);
    if (rc!=BLS_OK)
    {
        return rc;
    }

#ifdef DEBUG
    printf("pqnist_bls_sign SIG: ");
    OCT_output(&SIG);
#endif

    return 0;
}

//  Verify a signature using the BLS algorithm
int pqnist_bls_verify(char* M, char* pk, char* S)
{
    octet SIG = {SIGLEN,SIGLEN,S};
    octet PK = {G2LEN,G2LEN,pk};

#ifdef DEBUG
    printf("pqnist_bls_verify M %s\n", M);
    printf("pqnist_bls_verify PK: ");
    OCT_output(&PK);
    printf("pqnist_bls_verify SIG: ");
    OCT_output(&SIG);
#endif

    int rc=BLS_BLS381_VERIFY(&SIG,M,&PK);
    if (rc!=BLS_OK)
    {
        return rc;
    }

    return 0;
}

//  Add two members from the group G1
int pqnist_bls_addg1(char* r1, char* r2, char* r)
{
    octet R1 = {BGS_BLS381,BGS_BLS381,r1};
    octet R2 = {BGS_BLS381,BGS_BLS381,r2};
    octet R = {BGS_BLS381,BGS_BLS381,r};

    int rc=BLS_BLS381_ADD_G1(&R1,&R2,&R);
    if (rc!=BLS_OK)
    {
        return rc;
    }

#ifdef DEBUG
    printf("pqnist_bls_addg1 R1: ");
    OCT_output(&R1);
    printf("pqnist_bls_addg1 R2: ");
    OCT_output(&R2);
    printf("pqnist_bls_addg1 R: ");
    OCT_output(&R);
#endif

    return 0;
}

//  Add two members from the group G2
int pqnist_bls_addg2(char* r1, char* r2, char* r)
{
    octet R1 = {G2LEN,G2LEN,r1};
    octet R2 = {G2LEN,G2LEN,r2};
    octet R = {G2LEN,G2LEN,r};

    int rc=BLS_BLS381_ADD_G2(&R1,&R2,&R);
    if (rc!=BLS_OK)
    {
        return rc;
    }

#ifdef DEBUG
    printf("pqnist_bls_addg2 R1: ");
    OCT_output(&R1);
    printf("pqnist_bls_addg2 R2: ");
    OCT_output(&R2);
    printf("pqnist_bls_addg2 R: ");
    OCT_output(&R);
#endif

    return 0;
}

// Use Shamir's secret sharing to distribute BLS secret keys
int pqnist_bls_make_shares(int k, int n,  char* pSEED, char* pX, char* pY, char* pSKI, char* pSKO)
{
    int rc;

    octet SEED = {PQNIST_SEED_LENGTH,PQNIST_SEED_LENGTH,pSEED};
    octet SKI = {BGS_BLS381,BGS_BLS381,pSKI};
    octet SKO = {BGS_BLS381,BGS_BLS381,pSKO};
    octet X[n];
    octet Y[n];
    for(int i=0; i<n; i++)
    {
        Y[i].max = BGS_BLS381;
        Y[i].len = BGS_BLS381;
        Y[i].val = &pY[i*BGS_BLS381];
        X[i].max = BGS_BLS381;
        X[i].len = BGS_BLS381;
        X[i].val = &pX[i*BGS_BLS381];
    }

    // CSPRNG
    csprng RNG;

    // initialise strong RNG
    CREATE_CSPRNG(&RNG,&SEED);

    // Make shares of BLS secret key
    rc = BLS_BLS381_MAKE_SHARES(k, n, &RNG, X, Y, &SKI, &SKO);
    if (rc)
    {
        return rc;
    }

#ifdef DEBUG
    printf("pqnist_keys SEED: ");
    OCT_output(&SEED);
    printf("\n");

    for(int i=0; i<n; i++)
    {
        printf("X[%d] ", i);
        OCT_output(&X[i]);
        printf("Y[%d] ", i);
        OCT_output(&Y[i]);
        printf("\n");
    }

    printf("SKI: ");
    OCT_output(&SKI);
    printf("SKO: ");
    OCT_output(&SKO);
#endif

    return 0;
}

// Use Shamir's secret sharing to recover a BLS secret key
int pqnist_bls_recover_secret(int k, char* pX, char* pY, char* pSK)
{
    int rc;

    octet SK = {BGS_BLS381,BGS_BLS381,pSK};
    octet X[k];
    octet Y[k];
    for(int i=0; i<k; i++)
    {
        Y[i].max = BGS_BLS381;
        Y[i].len = BGS_BLS381;
        Y[i].val = &pY[i*BGS_BLS381];
        X[i].max = BGS_BLS381;
        X[i].len = BGS_BLS381;
        X[i].val = &pX[i*BGS_BLS381];
    }

    // Recover BLS secret key
    rc = BLS_BLS381_RECOVER_SECRET(k, X, Y, &SK);
    if (rc)
    {
        return rc;
    }

#ifdef DEBUG
    printf("SK: ");
    OCT_output(&SK);
#endif

    return 0;
}

// Use Shamir's secret sharing to recover a BLS signature
int pqnist_bls_recover_signature(int k, char* pX, char* pY, char* pSIG)
{
    int rc;

    octet SIG = {SIGLEN,SIGLEN,pSIG};

    octet X[k];
    octet Y[k];
    for(int i=0; i<k; i++)
    {
        Y[i].max = SIGLEN;
        Y[i].len = SIGLEN;
        Y[i].val = &pY[(SIGLEN)*i];
        X[i].max = BGS_BLS381;
        X[i].len = BGS_BLS381;
        X[i].val = &pX[BGS_BLS381*i];
    }

    // Recover BLS signature
    rc = BLS_BLS381_RECOVER_SIGNATURE(k, X, Y, &SIG);
    if (rc)
    {
        return rc;
    }

#ifdef DEBUG
    printf("pqnist_bls_recover_signature SIG: ");
    OCT_output(&SIG);
#endif

    return 0;
}


// Generate SIKE and BLS public and private key pairs
int pqnist_sike_keys(char* seed, char* SIKEpk, char* SIKEsk)
{
    int rc;

    // Initialise KAT RNG
    rc = OQS_randombytes_switch_algorithm(OQS_RAND_alg_nist_kat);
    if (rc != OQS_SUCCESS)
    {
        return rc;
    }
    OQS_randombytes_nist_kat_init(seed, NULL, 256);

    // Generate SIKE key pair
    rc = OQS_KEM_sike_p751_keypair(SIKEpk, SIKEsk);
    if (rc != OQS_SUCCESS)
    {
        return rc;
    }


#ifdef DEBUG
    int i = OQS_KEM_sike_p751_length_public_key;
    printf("pqnist_keys SIKEpklen %d SIKEpk: ", i);
    amcl_print_hex(SIKEpk, i);
    i = OQS_KEM_sike_p751_length_secret_key;
    printf("pqnist_keys SIKE sklen %d SIKEsk: ", i);
    amcl_print_hex(SIKEsk, i);
    printf("\n");
#endif

    return 0;
}


/*   The  message is encrypted using AES-256. The key
     is generated inside this function as an output
     from the encapsulation function.
*/
int pqnist_encapsulate_encrypt(char* P, int Plen, char* IV, char* pk, char* ek)
{
    // AES-256 key
    uint8_t K[OQS_KEM_sike_p751_length_shared_secret];

#ifdef DEBUG
    printf("Plaintext %d P: ", Plen);
    amcl_print_hex(P, Plen);
    int i = OQS_KEM_sike_p751_length_public_key;
    printf("pklen %d pk: ", i);
    amcl_print_hex(pk, i);
#endif

    OQS_STATUS rc = OQS_KEM_sike_p751_encaps(ek, K, pk);
    if (rc != OQS_SUCCESS)
    {
        OQS_MEM_cleanse(K, OQS_KEM_sike_p751_length_shared_secret);
        return rc;
    }

#ifdef DEBUG
    i = OQS_KEM_sike_p751_length_ciphertext;
    printf("ek1 %d ek1: ", i);
    amcl_print_hex(ek, i);
    i = OQS_KEM_sike_p751_length_shared_secret;
    printf("K %d K: ", i);
    amcl_print_hex(K, i);
#endif

    // Encrypt plaintext
    pqnist_aes_cbc_encrypt(K, PQNIST_AES_KEY_LENGTH, IV, P, Plen);

#ifdef DEBUG
    printf("K: ");
    amcl_print_hex(K, PQNIST_AES_KEY_LENGTH);
    printf("IV: ");
    amcl_print_hex((uint8_t*)IV, PQNIST_AES_IV_LENGTH);
    // Ciphertext
    printf("C: ");
    amcl_print_hex((uint8_t*)P, Plen);
#endif

    return 0;
}

/*   Decapsulate the AES key and use it to decrypt the
     ciphertext. The plaintext is returned using the C
     parameter.
*/
int pqnist_decapsulate_decrypt(char* C, int Clen, char* IV, char* sk, char* ek)
{
    char sec[OQS_KEM_sike_p751_length_shared_secret*2];

    // Encapsulated secret is 24 byte therefore needs to be run twice to
    // generate 32 byte AES Key
    OQS_STATUS rc = OQS_KEM_sike_p751_decaps(sec, ek, sk);
    if (rc != OQS_SUCCESS)
    {
        OQS_MEM_cleanse(sec, OQS_KEM_sike_p751_length_secret_key);
        return rc;
    }

#ifdef DEBUG
    int i = OQS_KEM_sike_p751_length_shared_secret;
    printf("sec1 %d sec1: ", i);
    amcl_print_hex(sec, i);
    printf("sec2 %d sec2: ", i);
    amcl_print_hex(&sec[OQS_KEM_sike_p751_length_shared_secret], i);
#endif

    // Decrypt the ciphertext
    pqnist_aes_cbc_decrypt(sec, PQNIST_AES_KEY_LENGTH, IV, C, Clen);

    return 0;
}

//  AES encryption using GCM mode
void pqnist_aes_gcm_encrypt(char* K, int Klen, char* IV, int IVlen, char* A, int Alen, char* P, int Plen, char* C, char* T)
{
    gcm g;
    GCM_init(&g,Klen,K,IVlen,IV);
    GCM_add_header(&g,A,Alen);
    GCM_add_plain(&g,C,P,Plen);
    GCM_finish(&g,T);
}

// AES Decryption using GCM mode
void pqnist_aes_gcm_decrypt(char* K, int Klen, char* IV, int IVlen, char* A, int Alen, char* C, int Clen, char* P, char* T)
{
    gcm g;
    GCM_init(&g,Klen,K,IVlen,IV);
    GCM_add_header(&g,A,Alen);
    GCM_add_cipher(&g,P,C,Clen);
    GCM_finish(&g,T);
}

//   AES encryption using CBC mode
void pqnist_aes_cbc_encrypt(char* K, int Klen, char* IV, char* P, int Plen)
{
#ifdef DEBUG
    printf("pqnist_aes_cbc_encrypt Klen %d K: \n", Klen);
    amcl_print_hex(K, Klen);
    printf("pqnist_aes_cbc_encrypt IVlen %d IV: \n", PQNIST_AES_IV_LENGTH);
    amcl_print_hex(IV, PQNIST_AES_IV_LENGTH);
    printf("pqnist_aes_cbc_encrypt Plen %d P: \n", Plen);
    amcl_print_hex(P, Plen);
#endif

    int blockSize=16;
    amcl_aes a;
    AES_init(&a,CBC,Klen,K,IV);
    for (int i=0; i<(Plen/blockSize); i++)
    {
        AES_encrypt(&a,&P[i*blockSize]);
    }

#ifdef DEBUG
    printf("pqnist_aes_cbc_encrypt Clen %d C: \n", Plen);
    amcl_print_hex(P, Plen);
#endif

}

//   AES decryption using CBC mode
void pqnist_aes_cbc_decrypt(char* K, int Klen, char* IV, char* C, int Clen)
{
#ifdef DEBUG
    printf("pqnist_aes_cbc_decrypt Klen %d K: \n", Klen);
    amcl_print_hex(K, Klen);
    printf("pqnist_aes_cbc_decrypt IVlen %d IV: \n", PQNIST_AES_IV_LENGTH);
    amcl_print_hex(IV, PQNIST_AES_IV_LENGTH);
    printf("pqnist_aes_cbc_decrypt Clen %d C: \n", Clen);
    amcl_print_hex(C, Clen);
#endif

    int blockSize=16;
    amcl_aes a;
    AES_init(&a,CBC,Klen,K,IV);
    for (int i=0; i<(Clen/blockSize); i++)
    {
        AES_decrypt(&a,&C[i*blockSize]);
    }

#ifdef DEBUG
    printf("pqnist_aes_cbc_decrypt Plen %d P: \n", Clen);
    amcl_print_hex(C, Clen);
#endif
}
