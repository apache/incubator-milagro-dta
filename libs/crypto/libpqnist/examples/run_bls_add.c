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
   Sign a message and verify the signature. 
   Generate a public key from externally generated secret key
   Add public keys and signatures.
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
    int rc;

    // Message to be signed
    char message[] = "test message";

    // Seed values for CSPRNG
    char seed1[PQNIST_SEED_LENGTH];
    char seed2[PQNIST_SEED_LENGTH];

    // seed values
    char* seed1Hex = "3370f613c4fe81130b846483c99c032c17dcc1904806cc719ed824351c87b0485c05089aa34ba1e1c6bfb6d72269b150";
    char* seed2Hex = "46389f32b7cdebbbc46b7165d8fae888c9de444898390a939977e1a066256a6f465e7d76307178aef81ae0c6841f9b7c";
    amcl_hex2bin(seed1Hex, seed1, PQNIST_SEED_LENGTH*2);
    amcl_hex2bin(seed2Hex, seed2, PQNIST_SEED_LENGTH*2);
    printf("seed1: ");
    amcl_print_hex(seed1, PQNIST_SEED_LENGTH);
    printf("seed2: ");
    amcl_print_hex(seed2, PQNIST_SEED_LENGTH);
    printf("\n");

    // BLS keys
    char sk1[BGS_BLS381];
    char pktmp[G2LEN];
    char pk1[G2LEN];    
    char sk2[BGS_BLS381];
    char pk2[G2LEN];
    char pk12[G2LEN];

    // BLS signature
    char sig1[SIGLEN];
    char sig2[SIGLEN];
    char sig12[SIGLEN];

    rc = pqnist_bls_keys(seed1, pktmp, sk1);
    if (rc)
    {
        fprintf(stderr, "ERROR pqnist_keys rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    rc = pqnist_bls_keys(seed2, pk2, sk2);
    if (rc)
    {
        fprintf(stderr, "ERROR pqnist_keys rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    // Generate pk from sk 
    rc = pqnist_bls_keys(NULL, pk1, sk1);
    if (rc)
    {
        fprintf(stderr, "ERROR pqnist_keys rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    printf("sk1: ");
    amcl_print_hex(sk1, BGS_BLS381);    
    printf("pktmp: ");
    amcl_print_hex(pktmp, G2LEN);
    printf("pk1: ");
    amcl_print_hex(pk1, G2LEN);
    printf("sk2: ");
    amcl_print_hex(sk2, BGS_BLS381);    
    printf("pk2: ");
    amcl_print_hex(pk2, G2LEN);
    printf("\n");

    // Sign message
    rc = pqnist_bls_sign(message, sk1, sig1);
    if(rc != BLS_OK)
    {
        fprintf(stderr, "ERROR pqnist_bls_sign rc: %d\n", rc);
        printf("FAILURE\n");
        exit(EXIT_FAILURE);
    }

    rc = pqnist_bls_sign(message, sk2, sig2);
    if(rc != BLS_OK)
    {
        fprintf(stderr, "ERROR pqnist_bls_sign rc: %d\n", rc);
        printf("FAILURE\n");
        exit(EXIT_FAILURE);
    }

    printf("sig1: ");
    amcl_print_hex(sig1, SIGLEN);
    printf("sig2: ");
    amcl_print_hex(sig2, SIGLEN);
    printf("\n");

    // Verify message
    rc = pqnist_bls_verify(message, pk1, sig1);
    if (rc != BLS_OK)
    {
        fprintf(stderr, "ERROR pqnist_bls_verify rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    rc = pqnist_bls_verify(message, pk2, sig2);
    if (rc != BLS_OK)
    {
        fprintf(stderr, "ERROR pqnist_bls_verify rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    rc = pqnist_bls_addg1(sig1, sig2, sig12);
    if (rc != BLS_OK)
    {
        fprintf(stderr, "ERROR pqnist_bls_addg1 rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    rc = pqnist_bls_addg2(pk1, pk2, pk12);
    if (rc != BLS_OK)
    {
        fprintf(stderr, "ERROR pqnist_bls_addg1 rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    printf("pk12: ");
    amcl_print_hex(pk12, G2LEN);
    printf("sig12: ");
    amcl_print_hex(sig12, SIGLEN);
    printf("\n");

    rc = pqnist_bls_verify(message, pk12, sig12);
    if (rc != BLS_OK)
    {
        fprintf(stderr, "ERROR pqnist_bls_verify rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    printf("TEST PASSED\n");
    exit(EXIT_SUCCESS);
}
