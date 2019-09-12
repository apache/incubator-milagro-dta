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

/**
 * @file pqnist.h
 * @author Kealan McCusker
 * @brief crypto function declarations
 *
 */


#ifndef PQNIST_H
#define PQNIST_H

#define PQNIST_AES_KEY_LENGTH 32  //!< AES-256 key length
#define PQNIST_AES_IV_LENGTH 16   //!< AES-256 initialization vector length
#define PQNIST_SEED_LENGTH 48   //!< CSPRNG seed length

#ifdef __cplusplus
extern "C" {
#endif

/**  @brief Generate BLS keys

     Generate BLS public and private key

     @param seed             seed value for CSPRNG
     @param BLSpk            BLS public key
     @param BLSsk            BLS secret key. Generated externally if seed set to NULL
     @return                 Zero for success or else an error code
 */
int pqnist_bls_keys(char* seed, char* BLSpk, char* BLSsk);


/**  @brief Sign a message

     The message is signed using the BLS algorithm

     @param M            Message to be signed
     @param sk           BLS secret key
     @param S            Signature
     @return             Zero for success or else an error code
 */
int pqnist_bls_sign(char* M, char* sk, char* S);

/**  @brief Verify a signature

     Verify a signature using the BLS algorithm

     @param M            Message that was signed
     @param pk           BLS public key
     @param S            Signature
     @return             Zero for success or else an error code
 */
int pqnist_bls_verify(char* M, char* pk, char* S);

/**	@brief Add two members from the group G1
 *
	@param  r1  member of G1
	@param  r2  member of G1
	@param  r   member of G1. r = r1+r2
	@return     Zero for success or else an error code
 */
int pqnist_bls_addg1(char* r1, char* r2, char* r);

/**	@brief Add two members from the group G2
 *
	@param  r1  member of G2
	@param  r2  member of G2
	@param  r   member of G2. r = r1+r2
	@return     Zero for success or else an error code
 */
int pqnist_bls_addg2(char* r1, char* r2, char* r);


/**	@brief Use Shamir's secret sharing to distribute BLS secret keys
 *
	@param  k       Threshold
	@param  n       Number of shares
        @param  pSEED   seed value for CSPRNG - 48 bytes
	@param  pX      X values
	@param  pY      Y values. Valid BLS secret keys
	@param  pSKI    Input secret key to be shared. Ignored if set to NULL
	@param  pSKO    Secret key that is shared
	@return         Zero for success or else an error code
 */
int pqnist_bls_make_shares(int k, int n,  char* pSEED, char* pX, char* pY, char* pSKI, char* pSKO);

/**	@brief Use Shamir's secret sharing to recover a BLS secret key
 *
	@param  k    Threshold
	@param  pX   X values
	@param  pY   Y values. Valid BLS secret keys
	@param  pSK  Secret key that is recovered
	@return      Zero for success or else an error code
 */
int pqnist_bls_recover_secret(int k, char* pX, char* pY, char* pSK);

/**	@brief Use Shamir's secret sharing to recover a BLS signature
 *
	@param  k     Threshold
	@param  pX    X values
	@param  pY    Y values. Valid BLS signatures
	@param  pSIG  Signature that is recovered
	@return       Zero for success or else an error code
 */
int pqnist_bls_recover_signature(int k, char* pX, char* pY, char* pSIG);

/**  @brief AES-GCM Encryption

     AES encryption using GCM mode

     @param K            Key
     @param Klen         Key length in bytes
     @param IV           Initialization vector IV
     @param IVlen        IV length in bytes
     @param A            Additional authenticated data (AAD)
     @param Alen         AAD length in bytes
     @param P            Plaintext
     @param Plen         Plaintext length in bytes
     @param C            Ciphertext (same length as P)
     @param T            Authentication tag
 */
void pqnist_aes_gcm_encrypt(char* K, int Klen, char* IV, int IVlen, char* A, int Alen, char* P, int Plen, char* C, char* T);

/**  @brief AES-GCM Decryption

     AES decryption using GCM mode

     @param K            Key
     @param Klen         Key length in bytes
     @param IV           Initialization vector IV
     @param IVlen        IV length in bytes
     @param A            Additional authenticated data (AAD)
     @param Alen         AAD length in bytes
     @param C            Ciphertext
     @param Clen         Ciphertext length in bytes
     @param P            Plaintext  (same length as C)
     @param T            Authentication tag
 */
void pqnist_aes_gcm_decrypt(char* K, int Klen, char* IV, int IVlen, char* A, int Alen, char* C, int Clen, char* P, char* T);

/**  @brief AES-CBC Encryption

     AES encryption using CBC mode

     @param K            Key
     @param Klen         Key length in bytes
     @param IV           Initialization vector IV (16 bytes)
     @param P            Plaintext / Ciphertext must be a multiple of the block size (16)
     @param Plen         Plaintext length in bytes
 */
void pqnist_aes_cbc_encrypt(char* K, int Klen, char* IV, char* P, int Plen);

/**  @brief AES-CBC Decryption

     AES decryption using CBC mode

     @param K            Key
     @param Klen         Key length in bytes
     @param IV           Initialization vector IV (16 bytes)
     @param C            Ciphertext / Plaintext must be a multiple of the block size (16)
     @param Clen         Ciphertext length in bytes
 */
void pqnist_aes_cbc_decrypt(char* K, int Klen, char* IV, char* C, int Clen);

/**  @brief Generate SIKE keys

     Generate SIKE public and private key

     @param seed             seed value for CSPRNG - 48 bytes
     @param SIKEpk           SIKE public key
     @param SIKEsk           SIKE secret key
     @return                 Zero for success or else an error code
 */
int pqnist_sike_keys(char* seed, char* SIKEpk, char* SIKEsk);

/**  @brief Encrypt a message and encapsulate the AES Key for
     a recipient.

     The  message is encrypted using AES-256. The key
     is generated inside this function as an output
     from the encapsulation function. The ciphertext
     is returned using the P paramter.

     @param P            Plaintext to be encrypted / Ciphertext. Padded with zero.
     @param Plen         Plaintext length in bytes must be a multiple of the block size (16)
     @param IV           Initialization vector IV (16 bytes)
     @param pk           SIKE public key
     @param ek           Encapsulated key
     @return             Zero for success or else an error code
 */
int pqnist_encapsulate_encrypt(char* P, int Plen, char* IV, char* pk, char* ek);

/**  @brief Decapsulate the AES Key and decrypt the message

     Decapsulate the AES key and use it to decrypt the
     ciphertext. The plaintext is returned using the C
     parameter.

     @param C            Ciphertext to be decrypted / Plaintext
     @param Clen         Ciphertext length in bytes must be a multiple of the block size (16)
     @param IV           Initialization vector IV
     @param sk           SIKE secret key
     @param ek           Encapsulated key
     @return             Zero for success or else an error code
 */
int pqnist_decapsulate_decrypt(char* C, int Clen, char* IV, char* sk, char* ek);

#ifdef __cplusplus
}
#endif

#endif
