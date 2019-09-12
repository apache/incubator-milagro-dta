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
   BLS sign a message and verify the signature. Introduce errors.
*/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <amcl/utils.h>
#include <amcl/randapi.h>
#include <amcl/bls_BLS381.h>
#include <amcl/pqnist.h>

#define G2LEN 4*BFS_BLS381
#define SIGLEN BFS_BLS381+1

int main()
{
    int i,rc;

    // Seed value for CSPRNG
    char seed[PQNIST_SEED_LENGTH];

    // Message to be sent to Bob
    char p[] = "Hello Bob! This is a message from Alice";
    octet P = {0, sizeof(p), p};

    // non random seed value
    for (i=0; i<PQNIST_SEED_LENGTH; i++) seed[i]=i+1;
    printf("SEED: ");
    amcl_print_hex(seed, PQNIST_AES_KEY_LENGTH);
    printf("\n");

    // Generate BLS keys

    // Alice's BLS keys
    char BLSsk[BGS_BLS381];
    char BLSpk[G2LEN];

    rc = pqnist_bls_keys(seed, BLSpk, BLSsk);
    if (rc)
    {
        fprintf(stderr, "ERROR pqnist_keys rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    printf("BLS pklen %d pk: ", G2LEN);
    amcl_print_hex(BLSpk, G2LEN);
    printf("BLS sklen %d BLS sk: ", BGS_BLS381);
    amcl_print_hex(BLSsk, BGS_BLS381);
    printf("\n");

    // BLS signature
    char S[SIGLEN];

    // Alice signs message
    rc = pqnist_bls_sign(P.val, BLSsk, S);
    if(rc != BLS_OK)
    {
        fprintf(stderr, "ERROR pqnist_bls_sign rc: %d\n", rc);
        printf("FAILURE\n");
        exit(EXIT_FAILURE);
    }

    printf("Alice Slen %d SIG", SIGLEN);
    amcl_print_hex(S, SIGLEN);
    printf("\n");

    // Bob verifies message
    rc = pqnist_bls_verify(P.val, BLSpk, S);
    if (rc == BLS_OK)
    {
        printf("SUCCESS pqnist_bls_verify rc: %d\n", rc);
    }
    else
    {
        fprintf(stderr, "ERROR pqnist_bls_verify rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }


    printf("Bob P ");
    OCT_output(&P);
    printf("\n");

    // Bob verifies corrupted message. This should fail
    char tmp = P.val[0];
    P.val[0] = 5;
    rc = pqnist_bls_verify(P.val, BLSpk, S);
    if (rc == BLS_FAIL)
    {
        fprintf(stderr, "ERROR pqnist_bls_verify rc: %d\n", rc);
    }
    else
    {
        printf("SUCCESS pqnist_bls_verify rc: %d\n", rc);
        printf("TEST FAILED\n");
        exit(EXIT_FAILURE);
    }

    // Fix message
    P.val[0] = tmp;
    printf("Bob P ");
    OCT_output(&P);
    printf("\n");

    // Check signature is correct
    rc = pqnist_bls_verify(P.val, BLSpk, S);
    if (rc == BLS_OK)
    {
        printf("SUCCESS pqnist_bls_verify rc: %d\n", rc);
    }
    else
    {
        fprintf(stderr, "ERROR pqnist_bls_verify rc: %d\n", rc);
        printf("TEST FAILED\n");
        exit(EXIT_FAILURE);
    }

    // Bob verifies corrupted signature. This should fail
    S[0] = 0;
    rc = pqnist_bls_verify(P.val, BLSpk, S);
    if (rc == BLS_INVALID_G1)
    {

        fprintf(stderr, "ERROR pqnist_bls_verify rc: %d\n", rc);
    }
    else
    {
        printf("SUCCESS pqnist_bls_verify rc: %d\n", rc);
        printf("TEST FAILED\n");
        exit(EXIT_FAILURE);
    }

    // clear memory
    OCT_clear(&P);

    printf("TEST PASSED\n");
    exit(EXIT_SUCCESS);
}
