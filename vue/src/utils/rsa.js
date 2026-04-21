import { JSEncrypt } from 'jsencrypt'
import request from '@/utils/request'

export async function getPublicKey() {
  const res = await request.get('/api/auth/key')
  return { keyId: res.data.key_id, publicKey: res.data.pub_key }
}

export async function encryptWithServerKey(plainText) {
  const { keyId, publicKey } = await getPublicKey()
  const encryptor = new JSEncrypt()
  encryptor.setPublicKey(publicKey)
  const encrypted = encryptor.encrypt(plainText)
  if (!encrypted) {
    throw new Error('密码加密失败')
  }
  return { encrypted, keyId }
}

export function encryptText(publicKey, plainText) {
  const encryptor = new JSEncrypt()
  encryptor.setPublicKey(publicKey)
  const encrypted = encryptor.encrypt(plainText)
  if (!encrypted) {
    throw new Error('密码加密失败')
  }
  return encrypted
}
