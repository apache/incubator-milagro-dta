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

    // Message to be signed
    char message[] = "test message";

    // Seed value for CSPRNG
    char seed1[PQNIST_SEED_LENGTH];
    char seed2[PQNIST_SEED_LENGTH];

    // seed values
    char* seed1Hex = "3370f613c4fe81130b846483c99c032c17dcc1904806cc719ed824351c87b0485c05089aa34ba1e1c6bfb6d72269b150";
    char* seed2Hex = "46389f32b7cdebbbc46b7165d8fae888c9de444898390a939977e1a066256a6f465e7d76307178aef81ae0c6841f9b7c";
    amcl_hex2bin(seed1Hex, seed1, PQNIST_SEED_LENGTH*2);
    amcl_hex2bin(seed2Hex, seed2, PQNIST_SEED_LENGTH*2);
    printf("seed1: ");
    amcl_print_hex(seed1, PQNIST_SEED_LENGTH);
    printf("\n");
    printf("seed2: ");
    amcl_print_hex(seed2, PQNIST_SEED_LENGTH);
    printf("\n");

    // BLS keys
    char ski[BGS_BLS381];
    char sko[BGS_BLS381];
    char skr[BGS_BLS381];
    char pki[G2LEN];

    rc = pqnist_bls_keys(seed1, pki, ski);
    if (rc)
    {
        fprintf(stderr, "ERROR pqnist_keys rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    printf("pki: ");
    amcl_print_hex(pki, G2LEN);
    printf("ski: ");
    amcl_print_hex(ski, BGS_BLS381);
    printf("\n");

    // BLS signature
    char sigi[SIGLEN];

    // Sign message
    rc = pqnist_bls_sign(message, ski, sigi);
    if(rc != BLS_OK)
    {
        fprintf(stderr, "ERROR pqnist_bls_sign rc: %d\n", rc);
        printf("FAILURE\n");
        exit(EXIT_FAILURE);
    }

    printf("sigi ");
    amcl_print_hex(sigi, SIGLEN);
    printf("\n");

    // Verify signature
    rc = pqnist_bls_verify(message, pki, sigi);
    if (rc == BLS_OK)
    {
        printf("SUCCESS pqnist_bls_verify rc: %d\n", rc);
    }
    else
    {
        fprintf(stderr, "ERROR pqnist_bls_verify rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }

    // Secret shares
    char x[BGS_BLS381*n];
    char y[BGS_BLS381*n];

    // Make shares of BLS secret key
    rc = pqnist_bls_make_shares(k, n, seed2, x, y, ski, sko);
    if (rc!=BLS_OK)
    {
        fprintf(stderr, "ERROR pqnist_bls_make_shares rc: %d\n", rc);
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
        fprintf(stderr, "ERROR pqnist_bls_recover_secret rc: %d\n", rc);
        exit(EXIT_FAILURE);
    }
    printf("skr: ");
    amcl_print_hex(skr, BGS_BLS381);
    printf("\n");

    // Generate public keys and signatures using shares
    char sigs[(SIGLEN)*n];
    char sigr[SIGLEN];
    for(int i=0; i<n; i++)
    {
        rc = pqnist_bls_sign(message, &y[BFS_BLS381*i], &sigs[(SIGLEN)*i]);
        if(rc != BLS_OK)
        {
            fprintf(stderr, "ERROR pqnist_bls_sign rc: %d\n", rc);
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

    printf("TEST PASSED\n");
    exit(EXIT_SUCCESS);
}
