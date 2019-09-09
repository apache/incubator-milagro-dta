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
   BLS SSS example
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
    int n=4;
    int k=3;

    // Seed value for CSPRNG
    char seed[PQNIST_SEED_LENGTH];

    // Message to be sent to Bob
    char p[] = "This is a test message";
    octet P = {0, sizeof(p), p};

    // non random seed value
    for (int i=0; i<PQNIST_SEED_LENGTH; i++) seed[i]=i+1;
    printf("SEED: ");
    amcl_print_hex(seed, PQNIST_SEED_LENGTH);
    printf("\n");

    // BLS keys
    char ski[BGS_BLS381];
    char sko[BGS_BLS381];
    char skr[BGS_BLS381];
    char pki[G2LEN];

    // BLS signature
    char sigi[SIGLEN];
    char sigr[SIGLEN];

    // Octets for testing
    octet SKI = {BGS_BLS381, BGS_BLS381, ski};
    octet SKR = {BGS_BLS381, BGS_BLS381, skr};
    octet SIGI = {SIGLEN, SIGLEN, sigi};
    octet SIGR = {SIGLEN, SIGLEN, sigr};

    rc = pqnist_bls_keys(seed, pki, ski);
    if (rc)
    {
        fprintf(stderr, "FAILURE pqnist_keys rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    printf("pki: ");
    amcl_print_hex(pki, G2LEN);
    printf("ski: ");
    amcl_print_hex(ski, BGS_BLS381);
    printf("\n");

    // Alice signs message
    rc = pqnist_bls_sign(P.val, ski, sigi);
    if(rc != BLS_OK)
    {
        fprintf(stderr, "FAILURE pqnist_bls_sign rc: %d\n", rc);
        printf("FAILURE\n");
        exit(EXIT_FAILURE);
    }

    printf("sigi ");
    amcl_print_hex(sigi, SIGLEN);
    printf("\n");

    // Bob verifies message
    rc = pqnist_bls_verify(P.val, pki, sigi);
    if (rc != BLS_OK)
    {
        fprintf(stderr, "FAILURE pqnist_bls_verify rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    // Secret shares
    char x[BGS_BLS381*n];
    char y[BGS_BLS381*n];

    // Make shares of BLS secret key
    rc = pqnist_bls_make_shares(k, n, seed, x, y, ski, sko);
    if (rc!=BLS_OK)
    {
        fprintf(stderr, "FAILURE pqnist_bls_make_shares rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    printf("ski: ");
    amcl_print_hex(ski, BGS_BLS381);
    printf("sko: ");
    amcl_print_hex(sko, BGS_BLS381);
    printf("\n");

    for(int i=0; i<n; i++)
    {
        printf("x[%d] ", i);
        amcl_print_hex(&x[i*BGS_BLS381], BGS_BLS381);
        printf("y[%d] ", i);
        amcl_print_hex(&y[i*BGS_BLS381], BGS_BLS381);
        printf("\n");
    }

    // Recover BLS secret key
    rc = pqnist_bls_recover_secret(k, x, y, skr);
    if (rc!=BLS_OK)
    {
        fprintf(stderr, "FAILURE pqnist_bls_recover_secret rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }
    printf("skr: ");
    amcl_print_hex(skr, BGS_BLS381);
    printf("\n");

    rc = OCT_comp(&SKI,&SKR);
    if(!rc)
    {
        fprintf(stderr, "FAILURE SKI != SKR rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }


    // Generate public keys and signatures using shares
    char sigs[(SIGLEN)*n];
    for(int i=0; i<n; i++)
    {
        rc = pqnist_bls_sign(P.val, &y[BFS_BLS381*i], &sigs[(SIGLEN)*i]);
        if(rc != BLS_OK)
        {
            fprintf(stderr, "FAILURE pqnist_bls_sign rc: %d\n", rc);
            printf("FAILURE\n");
            exit(EXIT_FAILURE);
        }

    }

    for(int i=0; i<n; i++)
    {
        printf("sigs[%d] ", i);
        amcl_print_hex(&sigs[i*(SIGLEN)], SIGLEN);
        printf("\n");
    }

    // Recover BLS signature
    pqnist_bls_recover_signature(k, x, sigs, sigr);

    printf("sigr ");
    amcl_print_hex(sigr, SIGLEN);
    printf("\n");

    rc = OCT_comp(&SIGI,&SIGR);
    if(!rc)
    {
        fprintf(stderr, "FAILURE SIGI != SIGR rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    printf("SUCCESS\n");
    exit(EXIT_SUCCESS);
}
