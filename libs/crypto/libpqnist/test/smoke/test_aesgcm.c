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
   Test AES-256 encryption and decryption in GCM mode
*/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <amcl/utils.h>
#include <amcl/randapi.h>
#include <pqnist/pqnist.h>

#define Keylen 32
#define IVlen 12
#define Taglen 12

int main()
{
    int i;

    // Seed value for CSPRNG
    // Seed value for CSPRNG
    char seed[PQNIST_SEED_LENGTH];
    octet SEED = {sizeof(seed),sizeof(seed),seed};

    csprng RNG;

    // Alice's authentication Tag
    char t1[Taglen];
    octet T1= {sizeof(t1),sizeof(t1),t1};

    // Bob's authentication Tag
    char t2[Taglen];
    octet T2= {sizeof(t2),sizeof(t2),t2};

    // AES Key
    char k[Keylen];
    octet K= {0,sizeof(k),k};

    // Initialization vector
    char iv[IVlen];
    octet IV= {0,sizeof(iv),iv};

    // Ciphertext
    char c[256];
    octet C= {0,sizeof(c),c};

    // Recovered plaintext
    char p2[256];
    octet P2= {0,sizeof(p2),p2};

    // Message to be sent to Bob
    char p1[256];
    octet P1 = {0, sizeof(p1), p1};
    OCT_jstring(&P1,"Hello Bob!");

    // Additional authenticated data (AAD)
    char aad[256];
    octet AAD = {0, sizeof(aad), aad};
    OCT_jstring(&AAD,"Header info");

    // non random seed value
    for (i=0; i<PQNIST_SEED_LENGTH; i++) SEED.val[i]=i+1;
    printf("SEED: ");
    OCT_output(&SEED);
    printf("\n");

    // initialise random number generator
    CREATE_CSPRNG(&RNG,&SEED);

    // Generate 256 bit AES Key
    K.len=Keylen;
    generateRandom(&RNG,&K);

    // Alice

    printf("Alice Key: ");
    amcl_print_hex(K.val, K.len);

    // Random initialization value
    IV.len=IVlen;
    generateRandom(&RNG,&IV);
    printf("Alice IV: ");
    OCT_output(&IV);

    printf("Alice AAD: ");
    OCT_output(&AAD);

    printf("Alice Plaintext: ");
    OCT_output(&P1);

    printf("Alice Plaintext: ");
    OCT_output_string(&P1);
    printf("\n");

    // Encrypt plaintext
    pqnist_aes_gcm_encrypt(K.val, K.len, IV.val, IV.len, AAD.val, AAD.len, P1.val, P1.len, C.val, T1.val);

    C.len = P1.len;
    printf("Alice: Ciphertext: ");
    OCT_output(&C);

    T1.len = Taglen;
    printf("Alice Tag: ");
    OCT_output(&T1);
    printf("\n");

    // Bob

    printf("Bob Key: ");
    amcl_print_hex(K.val, K.len);

    printf("Bob IV ");
    OCT_output(&IV);

    printf("Bob AAD: ");
    OCT_output(&AAD);

    printf("Bob Ciphertext: ");
    OCT_output(&C);

    pqnist_aes_gcm_decrypt(K.val, K.len, IV.val, IVlen, AAD.val, AAD.len, C.val, C.len, P2.val, T2.val);

    printf("Bob Plaintext: ");
    P2.len = C.len;
    OCT_output(&P2);

    printf("Bob Plaintext: ");
    OCT_output_string(&P2);
    printf("\n");

    printf("Bob Tag: ");
    T2.len = Taglen;
    OCT_output(&T2);

    if (!OCT_comp(&P1,&P2))
    {
        printf("FAILURE Decryption");
        exit(EXIT_FAILURE);
    }

    if (!OCT_comp(&T1,&T2))
    {
        printf("FAILURE TAG mismatch");
        exit(EXIT_FAILURE);
    }

    /* clear memory */
    OCT_clear(&T1);
    OCT_clear(&T2);
    OCT_clear(&K);
    OCT_clear(&IV);
    OCT_clear(&C);
    OCT_clear(&P1);
    OCT_clear(&P2);

    KILL_CSPRNG(&RNG);

    printf("SUCCESS\n");
    exit(EXIT_SUCCESS);
}
