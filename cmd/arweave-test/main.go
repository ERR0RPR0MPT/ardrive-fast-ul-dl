package main

import (
	"fmt"
	"github.com/ERR0RPR0MPT/ardrive-fast-ul-dl/arweave_api"
)

func main() {
	//arNode := "https://arweave.net"
	//proxyUrl := "http://127.0.0.1:8888"
	//w, err := goar.NewWalletFromPath("./key/ardrive-wallet.json", arNode, proxyUrl) // your wallet private key
	//
	////anchor, err := w.Client.GetTransactionAnchor()
	//if err != nil {
	//	fmt.Println("error w.Client.GetTransactionAnchor:", err)
	//	return
	//}
	//
	//data, err := os.ReadFile("./build.bat")
	//if err != nil {
	//	fmt.Println("error os.ReadFile:", err)
	//	return
	//}
	//
	//reward, err := w.Client.GetTransactionPrice(len(data), nil)
	//if err != nil {
	//	fmt.Println("error w.Client.GetTransactionPrice:", err)
	//	return
	//}
	//
	//tags := []schema.Tag{
	//	{
	//		Name:  "testSendData",
	//		Value: "123",
	//	},
	//}
	//
	//tx := &schema.Transaction{
	//	Format:   2,
	//	Target:   "",
	//	Quantity: "0",
	//	Tags:     utils.TagsEncode(tags),
	//	Data:     utils.Base64Encode(data),
	//	DataSize: fmt.Sprintf("%d", len(data)),
	//	Reward:   fmt.Sprintf("%d", reward),
	//}
	//
	////tx.LastTx = anchor
	//tx.Owner = utils.Base64Encode(w.Signer.PubKey.N.Bytes())
	//
	//if err = utils.SignTransaction(tx, w.Signer.PrvKey); err != nil {
	//	fmt.Println("error utils.SignTransaction:", err)
	//	return
	//}
	//
	//fmt.Println("w.Signer.PubKey.N.Bytes():", utils.Base64Encode(w.Signer.PubKey.N.Bytes()))
	//fmt.Println("w.Signer.PrvKey.D.Bytes():", utils.Base64Encode(w.Signer.PrvKey.D.Bytes()))

	//// 将 []byte 转为十六进制字符串
	//hexString := hex.EncodeToString(w.Signer.PrvKey.D.Bytes())
	//
	//// 每 1 字节分隔一个空格
	//var spacedHexStrings []string
	//for i := 0; i < len(hexString); i += 2 {
	//	spacedHexStrings = append(spacedHexStrings, hexString[i:i+2])
	//}
	//result := strings.Join(spacedHexStrings, " ")
	//
	//fmt.Println("w.Signer.PrvKey.D.Bytes().Result:", result)

	//fmt.Println("tx.Signature:", tx.Signature)
	//fmt.Println("tx.ID:", tx.ID)

	//uploader, err := goar.CreateUploader(w.Client, tx, nil)
	//if err != nil {
	//	fmt.Println("error goar.CreateUploader:", err)
	//	return
	//}
	//
	//err = uploader.Once()
	//if err != nil {
	//	fmt.Println("error uploader.Once():", err)
	//	return
	//}

	//signer02, err := goar.NewSignerFromPath("./key/ardrive-wallet.json")
	//if err != nil {
	//	fmt.Println("error goar.NewSignerFromPath:", err)
	//	return
	//}
	//
	//b2, err := goar.NewBundler(signer02)
	//if err != nil {
	//	fmt.Println("error goar.NewBundler(signer02):", err)
	//	return
	//}
	//
	//item2, err := b2.CreateAndSignItem([]byte("ar foo"), "", "", []schema.Tag{{Name: "Content-Type", Value: "application/txt"}})
	//if err != nil {
	//	fmt.Println("error b2.CreateAndSignItem:", err)
	//	return
	//}
	//
	//fmt.Println(string(item2.Binary))

	//bin := []byte{0x01, 0x00, 0x25, 0xED, 0xB7, 0x88, 0xE0, 0x92, 0x37, 0x89, 0x8B, 0xE4, 0x45, 0xD3, 0xB9, 0x62, 0x42, 0x3C, 0xFE, 0x10, 0x3F, 0xC1, 0x5E, 0x0A, 0x2F, 0x01, 0x85, 0xCE, 0x47, 0x4B, 0xD3, 0x6C, 0xE4, 0x85, 0xEC, 0x94, 0x8D, 0xDF, 0xA4, 0xD7, 0x37, 0x0B, 0xCA, 0x99, 0xC4, 0xB5, 0x7F, 0xE5, 0x8E, 0x13, 0x19, 0xE8, 0xC5, 0xED, 0x7B, 0xEA, 0x39, 0x17, 0xCA, 0x4E, 0x7D, 0xD6, 0xEE, 0xB2, 0xF3, 0x38, 0xD2, 0x9A, 0x1B, 0x93, 0x15, 0xC5, 0x01, 0x06, 0xDF, 0xA0, 0x3E, 0x13, 0x95, 0x6E, 0x47, 0x0F, 0x8F, 0x4A, 0xFC, 0xCA, 0x32, 0xDD, 0xD8, 0x15, 0x6C, 0xD4, 0xF1, 0x3D, 0xCA, 0x9A, 0x00, 0x2A, 0x97, 0x96, 0xBD, 0xD2, 0x9E, 0x5C, 0x4A, 0x9D, 0xF7, 0x4A, 0x8C, 0x8C, 0xB7, 0xF4, 0xE5, 0xC9, 0x8B, 0x0F, 0x75, 0x0A, 0xDC, 0x7D, 0x66, 0x2A, 0x27, 0x12, 0x7B, 0x6F, 0x83, 0x42, 0x02, 0x86, 0x08, 0x78, 0x0C, 0xE2, 0x68, 0xFE, 0x20, 0x08, 0x1F, 0xC7, 0x2C, 0xD1, 0xF3, 0xF0, 0xF4, 0xFC, 0xC0, 0x55, 0x76, 0x44, 0x3F, 0x18, 0x1C, 0x06, 0xEE, 0x44, 0xFC, 0x08, 0x77, 0x79, 0xC9, 0x87, 0x42, 0x2D, 0xEB, 0x16, 0x21, 0x26, 0x3E, 0xB8, 0xFA, 0x12, 0xAE, 0xC8, 0x54, 0x7D, 0xA5, 0x04, 0x67, 0x65, 0x77, 0xD0, 0x48, 0xD4, 0x4B, 0x1A, 0x60, 0x5D, 0x36, 0x1C, 0x8A, 0x9B, 0x7A, 0xF1, 0x76, 0x2E, 0xF9, 0x3A, 0xF5, 0xC6, 0x3D, 0x24, 0x10, 0x30, 0x52, 0xDD, 0x96, 0xE4, 0xD4, 0x2A, 0x59, 0xF8, 0xBF, 0xB0, 0xEE, 0xC1, 0x16, 0x2E, 0xF0, 0x94, 0x2B, 0x38, 0xDE, 0x00, 0x74, 0x21, 0x61, 0x14, 0x91, 0xD1, 0xEE, 0x7A, 0xF3, 0xFF, 0x57, 0x9B, 0xF5, 0x8D, 0x31, 0xE6, 0xA4, 0xFE, 0x0F, 0x10, 0x33, 0x2D, 0x30, 0x53, 0xCB, 0x15, 0xC3, 0xF6, 0xA9, 0x17, 0x7E, 0x79, 0xF9, 0x5E, 0x03, 0x02, 0xB6, 0xBA, 0xE0, 0xF4, 0xF2, 0x00, 0x5A, 0x89, 0xAC, 0xC2, 0xC5, 0xDA, 0xBB, 0x02, 0x70, 0x92, 0x15, 0x7B, 0xB4, 0xC0, 0x70, 0x7E, 0x8B, 0x60, 0x83, 0x34, 0x3F, 0x27, 0x45, 0xCD, 0x2D, 0x80, 0x0C, 0x35, 0xEC, 0x95, 0x4F, 0xAD, 0x3C, 0xB9, 0x14, 0x0D, 0x93, 0xB3, 0x9B, 0xFC, 0xA4, 0x26, 0x20, 0x20, 0x1D, 0xEC, 0x2C, 0xE1, 0xCE, 0xEA, 0xDD, 0x68, 0xE5, 0x1B, 0xF8, 0x0C, 0xB7, 0x24, 0x74, 0x76, 0x05, 0x2A, 0xD9, 0xA0, 0x70, 0x80, 0xA3, 0x5B, 0x00, 0xEB, 0x26, 0x79, 0xAF, 0xD5, 0x8A, 0x83, 0xE0, 0xBA, 0xC5, 0x4D, 0x18, 0xB4, 0x64, 0x26, 0x19, 0xD1, 0xF0, 0xE8, 0x7F, 0x27, 0x8F, 0xDD, 0x00, 0x17, 0xD9, 0xEB, 0x3E, 0x70, 0x9C, 0x72, 0x22, 0xC3, 0xB2, 0xC0, 0x41, 0x52, 0x49, 0xB3, 0xD4, 0xC9, 0x81, 0x27, 0x57, 0x5B, 0xCE, 0x80, 0x74, 0x10, 0x0F, 0x80, 0x26, 0x4F, 0xEA, 0x15, 0xE5, 0x04, 0x91, 0xA0, 0x1E, 0x1F, 0xEC, 0x7E, 0x55, 0x40, 0x8A, 0x1B, 0xB1, 0x37, 0x2E, 0xCE, 0x82, 0xF0, 0x83, 0xF5, 0xF5, 0x27, 0xFA, 0x04, 0xA4, 0xD2, 0xF0, 0x9A, 0x13, 0xC9, 0x88, 0xA8, 0x84, 0x24, 0xBE, 0x77, 0x6D, 0x8C, 0x85, 0x45, 0xF2, 0xD5, 0xDC, 0x28, 0x36, 0x9F, 0xAE, 0xC8, 0xC6, 0xEE, 0xCD, 0x57, 0x76, 0xEF, 0xDB, 0xF1, 0x8B, 0xCF, 0x9B, 0x4F, 0xCB, 0x48, 0xE9, 0x60, 0x2C, 0xF3, 0x0D, 0x60, 0xF6, 0xF3, 0xF5, 0xB4, 0xC3, 0xBD, 0xD7, 0xB6, 0x24, 0xA6, 0x29, 0x4D, 0x7B, 0x82, 0xBB, 0x83, 0x32, 0x5B, 0xCC, 0x49, 0x7B, 0x1B, 0x59, 0xDA, 0xC8, 0xE4, 0x82, 0x14, 0x24, 0x0E, 0xAE, 0x3D, 0x7D, 0xE6, 0x0B, 0x31, 0xC5, 0xDA, 0xE7, 0x2F, 0x2F, 0x3F, 0x10, 0x7D, 0x5B, 0x53, 0xF1, 0xD6, 0x44, 0x2D, 0xE8, 0xB8, 0xE0, 0xD9, 0xF2, 0x41, 0x8A, 0xE1, 0x0E, 0xD7, 0xDC, 0x6D, 0x79, 0xA9, 0x09, 0x92, 0xDD, 0xD0, 0x2E, 0x8F, 0x50, 0x55, 0xE2, 0x31, 0x4B, 0x42, 0x09, 0x78, 0x9C, 0x88, 0xFA, 0xE0, 0xC0, 0x3E, 0x6A, 0x6C, 0xC4, 0x55, 0x6A, 0x3F, 0x8A, 0xD8, 0x50, 0x34, 0x7B, 0x4E, 0xDC, 0x73, 0x57, 0xDF, 0x13, 0xD4, 0x36, 0xF0, 0xB7, 0x90, 0x90, 0x6A, 0xE4, 0x57, 0x94, 0x0A, 0x46, 0xCD, 0x76, 0x3F, 0x27, 0x19, 0x05, 0xF8, 0x8A, 0x2D, 0x83, 0x60, 0xD4, 0x3F, 0x34, 0x17, 0x81, 0xBD, 0x10, 0x19, 0x2A, 0x28, 0x86, 0x1D, 0x26, 0x40, 0xA4, 0xE7, 0x72, 0x37, 0x4A, 0x10, 0x30, 0x87, 0xC5, 0xB7, 0xD1, 0x01, 0x89, 0x57, 0x7F, 0x3D, 0x00, 0x00, 0x8C, 0x7A, 0x3F, 0xEF, 0x6B, 0x38, 0x59, 0x77, 0xA6, 0x64, 0x82, 0x87, 0x12, 0x77, 0x28, 0x5F, 0xBA, 0x6A, 0x0E, 0xE2, 0xFE, 0x76, 0xAD, 0xE9, 0x4D, 0x4B, 0x48, 0x3B, 0x1E, 0x51, 0x51, 0x2D, 0xFF, 0x92, 0x62, 0xCE, 0x0F, 0x87, 0x78, 0x26, 0x5C, 0xE6, 0x44, 0x8E, 0xE0, 0xAC, 0x09, 0xA4, 0x56, 0x1B, 0x99, 0x96, 0xEC, 0x62, 0xA7, 0x13, 0x40, 0x0A, 0xB0, 0x63, 0xD2, 0x8F, 0xC2, 0x06, 0x6A, 0xDD, 0x96, 0xAA, 0xDB, 0x45, 0x6E, 0x3C, 0xCE, 0x4A, 0xD2, 0x45, 0xCE, 0x5C, 0xC5, 0x42, 0xA5, 0xAA, 0x16, 0x92, 0xB2, 0x06, 0xAF, 0xE7, 0x19, 0xD4, 0x84, 0xA1, 0xEB, 0x0A, 0x76, 0x6C, 0xC6, 0x01, 0x74, 0x5E, 0x8D, 0x2D, 0x02, 0x06, 0x40, 0x47, 0x00, 0x31, 0xCA, 0x41, 0x1A, 0x88, 0x52, 0x03, 0xF4, 0x88, 0x16, 0x49, 0x31, 0x9A, 0x62, 0x99, 0xEC, 0x28, 0x02, 0x35, 0x98, 0x7A, 0x28, 0x99, 0x3F, 0xA5, 0xF4, 0x42, 0x0F, 0xBB, 0x44, 0x12, 0x33, 0xE1, 0x40, 0xFD, 0x8F, 0x20, 0x33, 0x7F, 0x2C, 0xF4, 0xA5, 0x2E, 0x3C, 0x9F, 0x9F, 0x0F, 0xD2, 0xD3, 0xB7, 0x53, 0xFA, 0x1C, 0x0C, 0x68, 0xCC, 0x9C, 0x47, 0xC0, 0x6F, 0x98, 0xD2, 0x26, 0x0B, 0x8D, 0x81, 0x41, 0xF1, 0xA2, 0xD2, 0x08, 0x78, 0x27, 0xEA, 0x42, 0x9F, 0xBC, 0x8D, 0xD2, 0xAC, 0x2D, 0x12, 0x49, 0x7C, 0x3A, 0x62, 0x61, 0xCA, 0x1F, 0xFD, 0xBF, 0xD3, 0xCC, 0x4B, 0x41, 0x80, 0x5F, 0x7C, 0xE1, 0xD2, 0xC8, 0x60, 0xA9, 0xA9, 0x3D, 0x6C, 0xBC, 0xBE, 0xF1, 0xD7, 0xDF, 0xD3, 0xA0, 0x08, 0x52, 0x22, 0xCA, 0x64, 0x80, 0xF4, 0xCD, 0xA3, 0xF9, 0x9F, 0x5E, 0xBB, 0x80, 0x51, 0x75, 0xF3, 0x12, 0x8D, 0x06, 0xA6, 0xD1, 0x66, 0x51, 0xE6, 0x03, 0x50, 0x9E, 0xC3, 0xDA, 0xFB, 0x98, 0x85, 0xFE, 0x93, 0x76, 0x0D, 0x87, 0x47, 0xA5, 0x99, 0x77, 0xAB, 0xF0, 0x3E, 0x7A, 0x93, 0x9A, 0x42, 0x29, 0xF5, 0xEB, 0x23, 0xAB, 0xBC, 0x0B, 0x02, 0xE3, 0xAA, 0x11, 0xC4, 0xAD, 0x99, 0x99, 0x86, 0x9F, 0x60, 0x8D, 0xFA, 0xF3, 0xD5, 0x81, 0x2D, 0x02, 0x4B, 0x7C, 0x70, 0xB6, 0x0D, 0x4B, 0x6C, 0xD3, 0xC8, 0xDD, 0xA8, 0x61, 0x6E, 0x5C, 0xCE, 0x7E, 0xED, 0xFF, 0xB0, 0x24, 0x32, 0x06, 0x04, 0x68, 0xA0, 0x7D, 0x20, 0xAA, 0xDF, 0x3F, 0xE4, 0x23, 0x94, 0x85, 0x3F, 0xB7, 0x68, 0x04, 0xFD, 0xCE, 0x18, 0x18, 0x51, 0x78, 0xAB, 0x7B, 0xB4, 0x00, 0x08, 0x76, 0xDE, 0xC1, 0xFD, 0xF2, 0x72, 0x47, 0x71, 0x29, 0x7C, 0x54, 0x72, 0x53, 0xF5, 0x04, 0x98, 0x71, 0x5C, 0xBD, 0x6A, 0xFB, 0x18, 0xA3, 0xC7, 0x26, 0x5B, 0x93, 0xAF, 0x66, 0xDD, 0x83, 0x88, 0xEC, 0xE3, 0xAB, 0xF1, 0x03, 0x18, 0xBE, 0x21, 0xA0, 0x1F, 0xD2, 0x25, 0x16, 0x5D, 0x25, 0x19, 0x20, 0xC5, 0x6C, 0x2E, 0xDE, 0xF8, 0x27, 0x9C, 0xE3, 0x84, 0xB7, 0xA0, 0xC9, 0xEF, 0x20, 0x0B, 0x3C, 0xEF, 0xA5, 0xCE, 0x9F, 0x00, 0x00, 0x0A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1E, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0x18, 0x43, 0x6F, 0x6E, 0x74, 0x65, 0x6E, 0x74, 0x2D, 0x54, 0x79, 0x70, 0x65, 0x20, 0x61, 0x70, 0x70, 0x6C, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6F, 0x6E, 0x2F, 0x6A, 0x73, 0x6F, 0x6E, 0x08, 0x41, 0x72, 0x46, 0x53, 0x08, 0x30, 0x2E, 0x31, 0x34, 0x16, 0x45, 0x6E, 0x74, 0x69, 0x74, 0x79, 0x2D, 0x54, 0x79, 0x70, 0x65, 0x0C, 0x66, 0x6F, 0x6C, 0x64, 0x65, 0x72, 0x10, 0x44, 0x72, 0x69, 0x76, 0x65, 0x2D, 0x49, 0x64, 0x48, 0x32, 0x38, 0x64, 0x32, 0x62, 0x34, 0x37, 0x65, 0x2D, 0x33, 0x38, 0x62, 0x63, 0x2D, 0x34, 0x64, 0x65, 0x66, 0x2D, 0x62, 0x30, 0x33, 0x62, 0x2D, 0x39, 0x61, 0x30, 0x39, 0x39, 0x66, 0x35, 0x31, 0x34, 0x62, 0x32, 0x39, 0x12, 0x46, 0x6F, 0x6C, 0x64, 0x65, 0x72, 0x2D, 0x49, 0x64, 0x48, 0x36, 0x61, 0x38, 0x63, 0x30, 0x62, 0x64, 0x35, 0x2D, 0x30, 0x33, 0x32, 0x62, 0x2D, 0x34, 0x35, 0x62, 0x32, 0x2D, 0x61, 0x39, 0x34, 0x39, 0x2D, 0x39, 0x64, 0x66, 0x33, 0x31, 0x35, 0x65, 0x39, 0x32, 0x62, 0x63, 0x62, 0x20, 0x50, 0x61, 0x72, 0x65, 0x6E, 0x74, 0x2D, 0x46, 0x6F, 0x6C, 0x64, 0x65, 0x72, 0x2D, 0x49, 0x64, 0x48, 0x61, 0x33, 0x39, 0x62, 0x32, 0x62, 0x66, 0x31, 0x2D, 0x65, 0x31, 0x35, 0x64, 0x2D, 0x34, 0x62, 0x32, 0x38, 0x2D, 0x38, 0x63, 0x33, 0x39, 0x2D, 0x66, 0x30, 0x31, 0x66, 0x65, 0x65, 0x36, 0x31, 0x37, 0x37, 0x63, 0x30, 0x10, 0x41, 0x70, 0x70, 0x2D, 0x4E, 0x61, 0x6D, 0x65, 0x16, 0x41, 0x72, 0x44, 0x72, 0x69, 0x76, 0x65, 0x2D, 0x41, 0x70, 0x70, 0x18, 0x41, 0x70, 0x70, 0x2D, 0x50, 0x6C, 0x61, 0x74, 0x66, 0x6F, 0x72, 0x6D, 0x06, 0x57, 0x65, 0x62, 0x16, 0x41, 0x70, 0x70, 0x2D, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6F, 0x6E, 0x0C, 0x32, 0x2E, 0x36, 0x37, 0x2E, 0x34, 0x12, 0x55, 0x6E, 0x69, 0x78, 0x2D, 0x54, 0x69, 0x6D, 0x65, 0x14, 0x31, 0x37, 0x34, 0x33, 0x33, 0x33, 0x30, 0x39, 0x30, 0x36, 0x00, 0x7B, 0x22, 0x6E, 0x61, 0x6D, 0x65, 0x22, 0x3A, 0x22, 0x32, 0x22, 0x7D}
	////bin := []byte{0x01, 0x00, 0x12, 0x67, 0x8B, 0x5A, 0xD6, 0xB4, 0x38, 0xB7, 0x74, 0x24, 0x07, 0xBF, 0xEC, 0xEB, 0x5A, 0x6A, 0xF7, 0x0F, 0x91, 0x68, 0xC3, 0x4F, 0x31, 0x68, 0x00, 0x90, 0xE4, 0x1E, 0xCF, 0x7C, 0xE6, 0x70, 0x03, 0xD6, 0x05, 0x88, 0xAF, 0xE2, 0x6E, 0x64, 0x4E, 0xCA, 0xBD, 0x62, 0x93, 0x50, 0xCC, 0xEA, 0xEB, 0x3F, 0x88, 0xCF, 0x8C, 0xF5, 0x1A, 0x76, 0x3A, 0x51, 0x79, 0x07, 0x05, 0x09, 0xF4, 0x1F, 0x69, 0x30, 0xB5, 0x44, 0xCB, 0x90, 0x59, 0x0E, 0xB4, 0x4C, 0x62, 0x30, 0xFC, 0x0D, 0x09, 0xC6, 0xC3, 0xFB, 0xB0, 0x1D, 0xD0, 0xC0, 0x60, 0x20, 0x80, 0x46, 0xD4, 0x5F, 0x8F, 0xD4, 0xDF, 0xE2, 0xBC, 0x62, 0x46, 0x98, 0xD1, 0xF0, 0x14, 0x9A, 0x52, 0x68, 0xCA, 0xD5, 0x92, 0x58, 0x78, 0xE5, 0x99, 0xFB, 0x89, 0xE3, 0x7B, 0x3B, 0x3B, 0xA3, 0xBC, 0xF8, 0xEF, 0xA7, 0x82, 0x10, 0xE4, 0x8B, 0x7B, 0x80, 0x4C, 0x33, 0x44, 0xA9, 0x64, 0x2C, 0xB5, 0xD8, 0x5A, 0x72, 0x05, 0xFE, 0x5C, 0x32, 0x72, 0x48, 0xE4, 0xC2, 0x9E, 0x06, 0xEB, 0xA5, 0x1E, 0xB5, 0x05, 0x16, 0x14, 0xE2, 0xE7, 0x29, 0xDB, 0x2B, 0x5A, 0x0D, 0x96, 0x36, 0x64, 0x1B, 0x2C, 0xEB, 0x15, 0x96, 0xD9, 0xF9, 0xA4, 0x4E, 0xAB, 0x4D, 0x64, 0xE0, 0xCF, 0xB9, 0x2F, 0xEF, 0x35, 0x04, 0x4A, 0x4C, 0xFF, 0xCE, 0x2A, 0x38, 0xDF, 0xF1, 0xC6, 0x88, 0x96, 0x23, 0xD6, 0xE4, 0x16, 0x2D, 0x0B, 0x2F, 0xA1, 0xCC, 0x18, 0x31, 0xC6, 0xDF, 0x24, 0x09, 0x0E, 0x61, 0x7C, 0xA1, 0x05, 0xCB, 0x1F, 0x6A, 0x3F, 0x42, 0x4C, 0xC5, 0x9E, 0xB4, 0x2B, 0xD7, 0x2B, 0xBF, 0x1B, 0xF2, 0x30, 0x95, 0xF9, 0x46, 0xDC, 0x11, 0x08, 0xFE, 0xED, 0x1B, 0xAE, 0x8C, 0xD7, 0xF6, 0x06, 0xE5, 0xD5, 0x91, 0x78, 0xE2, 0x3E, 0xF5, 0x8E, 0x9C, 0xE6, 0x60, 0xC2, 0xCD, 0xA7, 0x36, 0xED, 0x3E, 0xD3, 0x6A, 0x69, 0x38, 0xB7, 0x01, 0xEA, 0xE5, 0xDE, 0xF3, 0x99, 0xD1, 0x68, 0xC5, 0x0B, 0x3C, 0x5A, 0x22, 0x9E, 0x0E, 0x54, 0xCB, 0xB9, 0xAB, 0xB2, 0x1F, 0xC6, 0xA7, 0x97, 0x0C, 0x08, 0x93, 0x99, 0xD2, 0x51, 0xF9, 0x59, 0xFE, 0xD5, 0x51, 0xB2, 0x6C, 0x8A, 0x03, 0x0B, 0xA8, 0x50, 0x93, 0x06, 0x9A, 0x97, 0x53, 0x28, 0xDF, 0x14, 0x47, 0xF2, 0x8E, 0x8A, 0x82, 0x75, 0xD4, 0xB5, 0xAA, 0x3A, 0xBB, 0x6B, 0x67, 0x00, 0xB7, 0xCC, 0xC4, 0x49, 0x9F, 0x99, 0x0F, 0xEF, 0xCD, 0xD1, 0x95, 0xD4, 0xA7, 0xF6, 0x41, 0x5C, 0x51, 0x37, 0x2D, 0xA1, 0xAB, 0x8B, 0x06, 0xDD, 0xF4, 0x8D, 0x04, 0xCF, 0x2B, 0x68, 0x4D, 0x55, 0x83, 0x30, 0x5D, 0xE5, 0xAC, 0xE8, 0xBF, 0xB1, 0xEA, 0x24, 0x94, 0x3E, 0xD2, 0x24, 0xEA, 0x2E, 0xE0, 0xC6, 0x7D, 0x85, 0xFF, 0xAC, 0x66, 0x14, 0x95, 0x8A, 0xEF, 0xEA, 0x12, 0x01, 0xF4, 0x8D, 0xA8, 0xD3, 0xFC, 0x75, 0x71, 0x3B, 0x79, 0x45, 0x5D, 0x19, 0xB1, 0x69, 0xBF, 0x7E, 0xBC, 0xAC, 0x96, 0x4C, 0x1C, 0x58, 0x01, 0xC6, 0x74, 0xAC, 0xF1, 0x7D, 0x3C, 0x03, 0xE2, 0x5F, 0x01, 0x4B, 0xC1, 0xB6, 0xE3, 0x48, 0x70, 0xF4, 0xCF, 0xBE, 0x44, 0xEB, 0x9F, 0x54, 0x35, 0x62, 0x7F, 0x11, 0x87, 0x69, 0x4D, 0x3B, 0x30, 0xA9, 0xEC, 0x46, 0xD2, 0xD7, 0x93, 0xD9, 0x49, 0x4A, 0x3B, 0x92, 0x1E, 0x0C, 0x41, 0x3C, 0x15, 0x29, 0x50, 0x90, 0x50, 0x50, 0xCC, 0x24, 0x79, 0x79, 0x49, 0xE8, 0x56, 0xBC, 0x52, 0x46, 0xE4, 0x1D, 0xB0, 0x3A, 0x7F, 0x23, 0x81, 0x0C, 0x6B, 0x5A, 0x0E, 0xA5, 0xD1, 0xD2, 0x5A, 0x30, 0x3B, 0x29, 0x36, 0x30, 0x9A, 0x5F, 0x65, 0xCF, 0xDE, 0x8D, 0x63, 0xE2, 0x27, 0xDD, 0x18, 0x8A, 0xE1, 0x0E, 0xD7, 0xDC, 0x6D, 0x79, 0xA9, 0x09, 0x92, 0xDD, 0xD0, 0x2E, 0x8F, 0x50, 0x55, 0xE2, 0x31, 0x4B, 0x42, 0x09, 0x78, 0x9C, 0x88, 0xFA, 0xE0, 0xC0, 0x3E, 0x6A, 0x6C, 0xC4, 0x55, 0x6A, 0x3F, 0x8A, 0xD8, 0x50, 0x34, 0x7B, 0x4E, 0xDC, 0x73, 0x57, 0xDF, 0x13, 0xD4, 0x36, 0xF0, 0xB7, 0x90, 0x90, 0x6A, 0xE4, 0x57, 0x94, 0x0A, 0x46, 0xCD, 0x76, 0x3F, 0x27, 0x19, 0x05, 0xF8, 0x8A, 0x2D, 0x83, 0x60, 0xD4, 0x3F, 0x34, 0x17, 0x81, 0xBD, 0x10, 0x19, 0x2A, 0x28, 0x86, 0x1D, 0x26, 0x40, 0xA4, 0xE7, 0x72, 0x37, 0x4A, 0x10, 0x30, 0x87, 0xC5, 0xB7, 0xD1, 0x01, 0x89, 0x57, 0x7F, 0x3D, 0x00, 0x00, 0x8C, 0x7A, 0x3F, 0xEF, 0x6B, 0x38, 0x59, 0x77, 0xA6, 0x64, 0x82, 0x87, 0x12, 0x77, 0x28, 0x5F, 0xBA, 0x6A, 0x0E, 0xE2, 0xFE, 0x76, 0xAD, 0xE9, 0x4D, 0x4B, 0x48, 0x3B, 0x1E, 0x51, 0x51, 0x2D, 0xFF, 0x92, 0x62, 0xCE, 0x0F, 0x87, 0x78, 0x26, 0x5C, 0xE6, 0x44, 0x8E, 0xE0, 0xAC, 0x09, 0xA4, 0x56, 0x1B, 0x99, 0x96, 0xEC, 0x62, 0xA7, 0x13, 0x40, 0x0A, 0xB0, 0x63, 0xD2, 0x8F, 0xC2, 0x06, 0x6A, 0xDD, 0x96, 0xAA, 0xDB, 0x45, 0x6E, 0x3C, 0xCE, 0x4A, 0xD2, 0x45, 0xCE, 0x5C, 0xC5, 0x42, 0xA5, 0xAA, 0x16, 0x92, 0xB2, 0x06, 0xAF, 0xE7, 0x19, 0xD4, 0x84, 0xA1, 0xEB, 0x0A, 0x76, 0x6C, 0xC6, 0x01, 0x74, 0x5E, 0x8D, 0x2D, 0x02, 0x06, 0x40, 0x47, 0x00, 0x31, 0xCA, 0x41, 0x1A, 0x88, 0x52, 0x03, 0xF4, 0x88, 0x16, 0x49, 0x31, 0x9A, 0x62, 0x99, 0xEC, 0x28, 0x02, 0x35, 0x98, 0x7A, 0x28, 0x99, 0x3F, 0xA5, 0xF4, 0x42, 0x0F, 0xBB, 0x44, 0x12, 0x33, 0xE1, 0x40, 0xFD, 0x8F, 0x20, 0x33, 0x7F, 0x2C, 0xF4, 0xA5, 0x2E, 0x3C, 0x9F, 0x9F, 0x0F, 0xD2, 0xD3, 0xB7, 0x53, 0xFA, 0x1C, 0x0C, 0x68, 0xCC, 0x9C, 0x47, 0xC0, 0x6F, 0x98, 0xD2, 0x26, 0x0B, 0x8D, 0x81, 0x41, 0xF1, 0xA2, 0xD2, 0x08, 0x78, 0x27, 0xEA, 0x42, 0x9F, 0xBC, 0x8D, 0xD2, 0xAC, 0x2D, 0x12, 0x49, 0x7C, 0x3A, 0x62, 0x61, 0xCA, 0x1F, 0xFD, 0xBF, 0xD3, 0xCC, 0x4B, 0x41, 0x80, 0x5F, 0x7C, 0xE1, 0xD2, 0xC8, 0x60, 0xA9, 0xA9, 0x3D, 0x6C, 0xBC, 0xBE, 0xF1, 0xD7, 0xDF, 0xD3, 0xA0, 0x08, 0x52, 0x22, 0xCA, 0x64, 0x80, 0xF4, 0xCD, 0xA3, 0xF9, 0x9F, 0x5E, 0xBB, 0x80, 0x51, 0x75, 0xF3, 0x12, 0x8D, 0x06, 0xA6, 0xD1, 0x66, 0x51, 0xE6, 0x03, 0x50, 0x9E, 0xC3, 0xDA, 0xFB, 0x98, 0x85, 0xFE, 0x93, 0x76, 0x0D, 0x87, 0x47, 0xA5, 0x99, 0x77, 0xAB, 0xF0, 0x3E, 0x7A, 0x93, 0x9A, 0x42, 0x29, 0xF5, 0xEB, 0x23, 0xAB, 0xBC, 0x0B, 0x02, 0xE3, 0xAA, 0x11, 0xC4, 0xAD, 0x99, 0x99, 0x86, 0x9F, 0x60, 0x8D, 0xFA, 0xF3, 0xD5, 0x81, 0x2D, 0x02, 0x4B, 0x7C, 0x70, 0xB6, 0x0D, 0x4B, 0x6C, 0xD3, 0xC8, 0xDD, 0xA8, 0x61, 0x6E, 0x5C, 0xCE, 0x7E, 0xED, 0xFF, 0xB0, 0x24, 0x32, 0x06, 0x04, 0x68, 0xA0, 0x7D, 0x20, 0xAA, 0xDF, 0x3F, 0xE4, 0x23, 0x94, 0x85, 0x3F, 0xB7, 0x68, 0x04, 0xFD, 0xCE, 0x18, 0x18, 0x51, 0x78, 0xAB, 0x7B, 0xB4, 0x00, 0x08, 0x76, 0xDE, 0xC1, 0xFD, 0xF2, 0x72, 0x47, 0x71, 0x29, 0x7C, 0x54, 0x72, 0x53, 0xF5, 0x04, 0x98, 0x71, 0x5C, 0xBD, 0x6A, 0xFB, 0x18, 0xA3, 0xC7, 0x26, 0x5B, 0x93, 0xAF, 0x66, 0xDD, 0x83, 0x88, 0xEC, 0xE3, 0xAB, 0xF1, 0x03, 0x18, 0xBE, 0x21, 0xA0, 0x1F, 0xD2, 0x25, 0x16, 0x5D, 0x25, 0x19, 0x20, 0xC5, 0x6C, 0x2E, 0xDE, 0xF8, 0x27, 0x9C, 0xE3, 0x84, 0xB7, 0xA0, 0xC9, 0xEF, 0x20, 0x0B, 0x3C, 0xEF, 0xA5, 0xCE, 0x9F, 0x00, 0x00, 0x0A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1E, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0x18, 0x43, 0x6F, 0x6E, 0x74, 0x65, 0x6E, 0x74, 0x2D, 0x54, 0x79, 0x70, 0x65, 0x20, 0x61, 0x70, 0x70, 0x6C, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6F, 0x6E, 0x2F, 0x6A, 0x73, 0x6F, 0x6E, 0x08, 0x41, 0x72, 0x46, 0x53, 0x08, 0x30, 0x2E, 0x31, 0x34, 0x16, 0x45, 0x6E, 0x74, 0x69, 0x74, 0x79, 0x2D, 0x54, 0x79, 0x70, 0x65, 0x0C, 0x66, 0x6F, 0x6C, 0x64, 0x65, 0x72, 0x10, 0x44, 0x72, 0x69, 0x76, 0x65, 0x2D, 0x49, 0x64, 0x48, 0x32, 0x38, 0x64, 0x32, 0x62, 0x34, 0x37, 0x65, 0x2D, 0x33, 0x38, 0x62, 0x63, 0x2D, 0x34, 0x64, 0x65, 0x66, 0x2D, 0x62, 0x30, 0x33, 0x62, 0x2D, 0x39, 0x61, 0x30, 0x39, 0x39, 0x66, 0x35, 0x31, 0x34, 0x62, 0x32, 0x39, 0x12, 0x46, 0x6F, 0x6C, 0x64, 0x65, 0x72, 0x2D, 0x49, 0x64, 0x48, 0x66, 0x34, 0x32, 0x63, 0x39, 0x64, 0x38, 0x65, 0x2D, 0x64, 0x36, 0x63, 0x30, 0x2D, 0x34, 0x32, 0x33, 0x32, 0x2D, 0x38, 0x31, 0x32, 0x66, 0x2D, 0x39, 0x63, 0x33, 0x63, 0x32, 0x30, 0x35, 0x61, 0x34, 0x65, 0x31, 0x65, 0x20, 0x50, 0x61, 0x72, 0x65, 0x6E, 0x74, 0x2D, 0x46, 0x6F, 0x6C, 0x64, 0x65, 0x72, 0x2D, 0x49, 0x64, 0x48, 0x61, 0x33, 0x39, 0x62, 0x32, 0x62, 0x66, 0x31, 0x2D, 0x65, 0x31, 0x35, 0x64, 0x2D, 0x34, 0x62, 0x32, 0x38, 0x2D, 0x38, 0x63, 0x33, 0x39, 0x2D, 0x66, 0x30, 0x31, 0x66, 0x65, 0x65, 0x36, 0x31, 0x37, 0x37, 0x63, 0x30, 0x10, 0x41, 0x70, 0x70, 0x2D, 0x4E, 0x61, 0x6D, 0x65, 0x16, 0x41, 0x72, 0x44, 0x72, 0x69, 0x76, 0x65, 0x2D, 0x41, 0x70, 0x70, 0x18, 0x41, 0x70, 0x70, 0x2D, 0x50, 0x6C, 0x61, 0x74, 0x66, 0x6F, 0x72, 0x6D, 0x06, 0x57, 0x65, 0x62, 0x16, 0x41, 0x70, 0x70, 0x2D, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6F, 0x6E, 0x0C, 0x32, 0x2E, 0x36, 0x37, 0x2E, 0x34, 0x12, 0x55, 0x6E, 0x69, 0x78, 0x2D, 0x54, 0x69, 0x6D, 0x65, 0x14, 0x31, 0x37, 0x34, 0x33, 0x33, 0x34, 0x32, 0x31, 0x38, 0x37, 0x00, 0x7B, 0x22, 0x6E, 0x61, 0x6D, 0x65, 0x22, 0x3A, 0x22, 0x33, 0x22, 0x7D}
	//
	////eyJuYW1lIjoiMyJ9
	////eyJuYW1lIjoiMiJ9
	//
	//resBundle, err := utils.DecodeBundleItem(bin)
	//if err != nil {
	//	fmt.Println("error b2.CreateAndSignItem:", err)
	//	return
	//}
	//
	//err = utils.VerifyBundleItem(resBundle)
	//if err != nil {
	//	fmt.Println("error utils.VerifyBundleItem(resBundle):", err)
	//	return
	//}
	//
	//resBundleNoBin := resBundle
	//resBundleNoBin.Binary = []byte{}
	//
	//data, err := json.MarshalIndent(resBundleNoBin, "", "  ")
	//if err != nil {
	//	fmt.Println("error json.MarshalIndent:", err)
	//	return
	//}
	//
	////输出原始对象
	//fmt.Println(string(data))

	//data, err := json.MarshalIndent(resBundleNoBin, "", "  ")
	//if err != nil {
	//	fmt.Println("error json.MarshalIndent:", err)
	//	return
	//}
	//
	////输出原始对象
	//fmt.Println(string(data))

	//name, err := arweave_api.ArDriveTurboCreateFolder(
	//	"test",
	//	"28d2b47e-38bc-4def-b03b-9a099f514b29",
	//	"a39b2bf1-e15d-4b28-8c39-f01fee6177c0",
	//	"./key/ardrive-wallet.json",
	//)
	//if err != nil {
	//	fmt.Println("error arweave_api.ArDriveTurboCreateFolder:", err)
	//	return
	//}
	//fmt.Println("Success:", name)

	txId, err := arweave_api.ArDriveTurboUploadFile(
		"./go.mod",
		"./key/ardrive-wallet.json",
	)
	if err != nil {
		fmt.Println("error arweave_api.ArDriveTurboUploadFile:", err)
		return
	}
	fmt.Println("Success: txId:", txId)
}
