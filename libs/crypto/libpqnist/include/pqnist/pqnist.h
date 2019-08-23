/**
 * @file pqnist.h
 * @author Kealan McCusker
 * @brief envelope crypto function declarations
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

/**  @brief Generate SIKE and BLS keys

     Generate SIKE and BLS key public and private key pairs

     @param seed             seed value for CSPRNG - 48 bytes
     @param SIKEpk           SIKE public key
     @param SIKEsk           SIKE secret key
     @param BLSpk            BLS public key
     @param BLSsk            BLS secret key
     @return                 Zero for success or else an error code
 */
int pqnist_keys(char* seed, char* SIKEpk, char* SIKEsk, char* BLSpk, char* BLSsk);

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

/**  @brief Sign a message

     The message is signed using the BLS algorithm

     @param M            Message to be signed
     @param sk           BLS secret key
     @param S            Signature
     @return             Zero for success or else an error code
 */
int pqnist_sign(char* M, char* sk, char* S);

/**  @brief Verify a signature

     Verify a signature using the BLS algorithm

     @param M            Message that was signed
     @param pk           BLS public key
     @param S            Signature
     @return             Zero for success or else an error code
 */
int pqnist_verify(char* M, char* pk, char* S);

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

#ifdef __cplusplus
}
#endif

#endif
