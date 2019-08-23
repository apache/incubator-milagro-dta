/*
   Sign a message and verify the signature
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
#define SIGLEN BFS_BLS381+1

int main()
{
    int i,rc;

    // Seed value for CSPRNG
    char seed[NTHREADS][PQNIST_SEED_LENGTH];

    // Message to be sent to Bob
    char p[NTHREADS][MAXSIZE];
    octet P[NTHREADS];

    // BLS signature
    char s[NTHREADS][SIGLEN];
    octet S[NTHREADS];

    // Initialise seed
    for(i=0; i<NTHREADS; i++)
    {
        for(int j=0; j<PQNIST_SEED_LENGTH; j++)
        {
            seed[i][j] = i;
        }
    }

    // Generate SIKE and BLS keys

    // Bob's SIKE keys (not used)
    uint8_t SIKEpk[NTHREADS][OQS_KEM_sike_p751_length_public_key];
    uint8_t SIKEsk[NTHREADS][OQS_KEM_sike_p751_length_secret_key];

    // Alice's BLS keys
    char BLSsk[NTHREADS][BGS_BLS381];
    char BLSpk[NTHREADS][G2LEN];

    #pragma omp parallel for
    for(i=0; i<NTHREADS; i++)
    {

        rc = pqnist_keys(seed[i], SIKEpk[i], SIKEsk[i], BLSpk[i], BLSsk[i]);
        if (rc)
        {
            fprintf(stderr, "FAILURE pqnist_keys rc: %d\n", rc);
            exit(EXIT_FAILURE);
        }

        printf("BLS pklen %d pk: ", G2LEN);
        amcl_print_hex(BLSpk[i], G2LEN);
        printf("BLS sklen %d BLS sk: ", BGS_BLS381);
        amcl_print_hex(BLSsk[i], BGS_BLS381);
        printf("\n");
    }

    // Alice

    for(i=0; i<NTHREADS; i++)
    {
        bzero(p[i],sizeof(p[i]));
        P[i].max = MAXSIZE;
        P[i].len = sprintf(p[i], "Hello Bob! This is a message from Alice %d", i);
        P[i].val = p[i];
        printf("Alice Plaintext: ");
        OCT_output_string(&P[i]);
        printf("\n");
    }

    for(i=0; i<NTHREADS; i++)
    {
        bzero(s[i],sizeof(s[i]));
        S[i].max = SIGLEN;
        S[i].len = SIGLEN;
        S[i].val = s[i];
    }

    #pragma omp parallel for
    for(i=0; i<NTHREADS; i++)
    {
        // Alice signs message
        rc = pqnist_sign(P[i].val, BLSsk[i], S[i].val);
        if(rc)
        {
            fprintf(stderr, "FAILURE pqnist_sign rc: %d\n", rc);
            exit(EXIT_FAILURE);
        }
        printf("Alice SIGlen %d  SIG", S[i].len);
        OCT_output(&S[i]);
        printf("\n");
    }

    #pragma omp parallel for
    for(i=0; i<NTHREADS; i++)
    {
        // Bob verifies message
        rc = pqnist_verify(P[i].val, BLSpk[i], S[i].val);
        if (rc)
        {
            fprintf(stderr, "FAILURE pqnist_verify rc: %d\n", rc);
            exit(EXIT_FAILURE);
        }
        else
        {
            printf("Test %d pqnist_verify rc: %d\n", i, rc);
            OCT_output_string(&P[i]);
            printf("\n");
        }
    }

    // clear memory
    for(i=0; i<NTHREADS; i++)
    {
        OQS_MEM_cleanse(SIKEsk[i], OQS_KEM_sike_p751_length_secret_key);
        OCT_clear(&P[i]);
        OCT_clear(&S[i]);
    }

    printf("SUCCESS\n");
    exit(EXIT_SUCCESS);
}
