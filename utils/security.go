package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// EncryptMessage cifra un mensaje usando AES-GCM con una clave derivada
func EncryptMessage(message string, key string) string {
	// Derivar clave de 32 bytes usando SHA-256
	hash := sha256.Sum256([]byte(key))

	// Crear bloque AES
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return ""
	}

	// Crear GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return ""
	}

	// Generar nonce aleatorio
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return ""
	}

	// Cifrar el mensaje
	ciphertext := gcm.Seal(nonce, nonce, []byte(message), nil)

	// Retornar en base64
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// DecryptMessage descifra un mensaje usando AES-GCM con una clave derivada
func DecryptMessage(cipherText string, key string) string {
	// Decodificar de base64
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return ""
	}

	// Derivar clave de 32 bytes usando SHA-256
	hash := sha256.Sum256([]byte(key))

	// Crear bloque AES
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return ""
	}

	// Crear GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return ""
	}

	// Verificar tamaño mínimo
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return ""
	}

	// Extraer nonce y ciphertext
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Descifrar
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return ""
	}

	return string(plaintext)
}

// GenerateToken genera un token simple basado en timestamp y clave
func GenerateToken(nodeID string, secretKey string) string {
	timestamp := time.Now().Unix()
	data := fmt.Sprintf("%s:%d", nodeID, timestamp)

	// Crear hash del token
	hash := sha256.Sum256([]byte(data + secretKey))
	token := fmt.Sprintf("%s.%s", data, hex.EncodeToString(hash[:16]))

	return base64.StdEncoding.EncodeToString([]byte(token))
}

// ValidateToken valida un token generado previamente
func ValidateToken(token string, secretKey string) (bool, string) {
	// Decodificar base64
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false, ""
	}

	// Dividir token en partes
	parts := strings.Split(string(decoded), ".")
	if len(parts) != 2 {
		return false, ""
	}

	data, receivedHash := parts[0], parts[1]

	// Recrear hash esperado
	hash := sha256.Sum256([]byte(data + secretKey))
	expectedHash := hex.EncodeToString(hash[:16])

	// Verificar hash
	if receivedHash != expectedHash {
		return false, ""
	}

	// Extraer nodeID del data
	dataParts := strings.Split(data, ":")
	if len(dataParts) != 2 {
		return false, ""
	}

	nodeID := dataParts[0]
	return true, nodeID
}

// ValidateTokenMiddleware middleware para validar tokens en peticiones HTTP
func ValidateTokenMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obtener token del header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Token requerido", http.StatusUnauthorized)
				return
			}

			// Verificar formato "Bearer <token>"
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Formato de token inválido", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Validar token
			valid, nodeID := ValidateToken(token, secretKey)
			if !valid {
				http.Error(w, "Token inválido", http.StatusUnauthorized)
				return
			}

			// Agregar nodeID al contexto de la petición
			r.Header.Set("X-Node-ID", nodeID)

			// Continuar con el siguiente handler
			next.ServeHTTP(w, r)
		})
	}
}

// HashPassword crea un hash simple de una contraseña (para demostración)
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// VerifyPassword verifica una contraseña contra su hash
func VerifyPassword(password, hash string) bool {
	return HashPassword(password) == hash
}

// SecureMessage estructura para mensajes seguros entre nodos
type SecureMessage struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

// CreateSecureMessage crea un mensaje seguro y firmado
func CreateSecureMessage(from, to, content, key string) *SecureMessage {
	msg := &SecureMessage{
		From:      from,
		To:        to,
		Content:   content,
		Timestamp: time.Now().Unix(),
	}

	// Crear firma del mensaje
	data := fmt.Sprintf("%s:%s:%s:%d", from, to, content, msg.Timestamp)
	hash := sha256.Sum256([]byte(data + key))
	msg.Signature = hex.EncodeToString(hash[:])

	return msg
}

// VerifySecureMessage verifica la integridad de un mensaje seguro
func VerifySecureMessage(msg *SecureMessage, key string) error {
	if msg == nil {
		return errors.New("mensaje nulo")
	}

	// Recrear firma esperada
	data := fmt.Sprintf("%s:%s:%s:%d", msg.From, msg.To, msg.Content, msg.Timestamp)
	hash := sha256.Sum256([]byte(data + key))
	expectedSignature := hex.EncodeToString(hash[:])

	// Verificar firma
	if msg.Signature != expectedSignature {
		return errors.New("firma inválida")
	}

	// Verificar timestamp (no más de 5 minutos de antigüedad)
	if time.Now().Unix()-msg.Timestamp > 300 {
		return errors.New("mensaje expirado")
	}

	return nil
}
