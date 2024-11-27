package hashtable

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"io"
	"os"
)

// Item представляет элемент хэш-таблицы
type Item struct {
	key   string
	value string
	next  *Item
}

// HashTable представляет хэш-таблицу
type HashTable struct {
	sizeArr int
	tabl    []*Item
}

// NewHashTable создает новую хэш-таблицу с заданным размером
func NewHashTable(size int) *HashTable {
	return &HashTable{
		sizeArr: size,
		tabl:    make([]*Item, size),
	}
}

// Hash вычисляет хэш для заданного ключа
func (ht *HashTable) Hash(itemKey string) int {
	h := fnv.New32a()
	h.Write([]byte(itemKey))
	return int(h.Sum32()) % ht.sizeArr
}

// IsFull проверяет, заполнена ли хэш-таблица
func (ht *HashTable) IsFull() bool {
	count := 0
	for i := 0; i < ht.sizeArr; i++ {
		if ht.tabl[i] != nil {
			count++
		}
	}
	return count >= ht.sizeArr
}

// AddHash добавляет элемент в хэш-таблицу
func (ht *HashTable) AddHash(key, value string) {
	index := ht.Hash(key)

	// Проверка на наличие уже такого ключа
	current := ht.tabl[index]
	for current != nil {
		if current.key == key {
			fmt.Printf("Ключ '%s' уже существует. Значение не добавлено.\n", key)
			return
		}
		current = current.next
	}

	// Проверка на есть ли место
	if ht.IsFull() {
		fmt.Println("Хэш-таблица переполнена. Невозможно добавить новый элемент.")
		return
	}

	// Добавление элемента
	newItem := &Item{key: key, value: value, next: ht.tabl[index]}
	ht.tabl[index] = newItem
}

// KeyItem получает значение по ключу
func (ht *HashTable) KeyItem(key string) {
	index := ht.Hash(key)
	current := ht.tabl[index]
	for current != nil {
		if current.key == key {
			fmt.Printf("key: %s value: %s\n", key, current.value)
			return
		}
		current = current.next
	}
	fmt.Println("Такого ключа нет.")
}

// DelValue удаляет элемент по ключу
func (ht *HashTable) DelValue(key string) {
	index := ht.Hash(key)
	var prev *Item
	current := ht.tabl[index]
	for current != nil {
		if current.key == key {
			if prev == nil {
				ht.tabl[index] = current.next
			} else {
				prev.next = current.next
			}
			return
		}
		prev = current
		current = current.next
	}
	fmt.Println("Такого ключа нет.")
}

// Print выводит содержимое хэш-таблицы
func (ht *HashTable) Print() {
	for i := 0; i < ht.sizeArr; i++ {
		current := ht.tabl[i]
		for current != nil {
			fmt.Printf("key: %s value: %s\n", current.key, current.value)
			current = current.next
		}
	}
}

// SerializeBinary сериализует хэш-таблицу в бинарный формат
func (ht *HashTable) SerializeBinary(name string) error {
	// Открываем файл для записи в бинарном режиме
	file, check := os.Create(name)
	if check != nil {
		return check
	}
	defer file.Close()

	// Проходим по всем элементам хэш-таблицы
	for i := 0; i < ht.sizeArr; i++ {
		current := ht.tabl[i]
		for current != nil {
			// Получаем длину ключа и значения
			lenK := uint32(len(current.key))
			lenV := uint32(len(current.value))

			// Записываем длину ключа в файл
			check := binary.Write(file, binary.LittleEndian, lenK)
			if check != nil {
				return check
			}
			// Записываем ключ в файл
			_, check = file.Write([]byte(current.key))
			if check != nil {
				return check
			}

			// Записываем длину значения в файл
			check = binary.Write(file, binary.LittleEndian, lenV)
			if check != nil {
				return check
			}
			// Записываем значение в файл
			_, check = file.Write([]byte(current.value))
			if check != nil {
				return check
			}

			// Переходим к следующему элементу в цепочке
			current = current.next
		}
	}

	return nil
}

// DeserializeBinary десериализует хэш-таблицу из бинарного формата
func (ht *HashTable) DeserializeBinary(filename string) error {
	file, check := os.Open(filename)
	if check != nil {
		return check
	}
	defer file.Close()

	ht.tabl = make([]*Item, ht.sizeArr)

	for {
		var lenK, lenV uint32
		// Читаем len ключ
		check := binary.Read(file, binary.LittleEndian, &lenK)
		if check == io.EOF {
			break
		}
		if check != nil {
			return check
		}
		// Читаем ключ
		keyBytes := make([]byte, lenK)
		_, check = file.Read(keyBytes)
		if check != nil {
			return check
		}
		// Читаем len значение
		check = binary.Read(file, binary.LittleEndian, &lenV)
		if check != nil {
			return check
		}
		// Читаем значение
		valueBytes := make([]byte, lenV)
		_, check = file.Read(valueBytes)
		if check != nil {
			return check
		}
		// Преобразуем байты в строки
		key := string(keyBytes)
		value := string(valueBytes)

		ht.AddHash(key, value)
	}

	return nil
}

// SerializeText сериализует хэш-таблицу в текстовый формат
func (ht *HashTable) SerializeText(filename string) error {
	file, check := os.Create(filename)
	if check != nil {
		return check
	}
	defer file.Close()

	for i := 0; i < ht.sizeArr; i++ {
		current := ht.tabl[i]
		for current != nil {
			_, check := fmt.Fprintf(file, "%s %s\n", current.key, current.value)
			if check != nil {
				return check
			}
			current = current.next
		}
	}

	return nil
}

// DeserializeText десериализует хэш-таблицу из текстового формата
func (ht *HashTable) DeserializeText(filename string) error {
	file, check := os.Open(filename)
	if check != nil {
		return check
	}
	defer file.Close()

	ht.tabl = make([]*Item, ht.sizeArr)
	// Используем сканер для чтения строк из файла
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var key, value string
		// Разбираем строку на ключ и значение
		_, check := fmt.Sscanf(line, "%s %s", &key, &value)
		if check != nil {
			return check
		}

		ht.AddHash(key, value)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
