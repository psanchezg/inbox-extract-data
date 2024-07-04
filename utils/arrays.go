package utils

func InsertAtBeginning(slice [][]interface{}, value []interface{}) [][]interface{} {
	// Crear un nuevo slice con espacio adicional para el nuevo valor
	newSlice := make([][]interface{}, len(slice)+1)
	// Colocar el nuevo valor en la primera posición
	newSlice[0] = value
	// Copiar los elementos del slice original al nuevo slice
	copy(newSlice[1:], slice)
	return newSlice
}
