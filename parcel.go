package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	// ParcelStatusRegistered indicates that the parcel has been registered.
	ParcelStatusRegistered = "registered"
	// ParcelStatusSent indicates that the parcel has been sent.
	ParcelStatusSent = "sent"
	// ParcelStatusDelivered indicates that the parcel has been delivered.
	ParcelStatusDelivered = "delivered"
)

// Parcel struct represents the information of a parcel.
type Parcel struct {
	// Number is a unique identifier for the parcel.
	Number int64 `json:"number"`
	// Client is the identifier of the client who ordered the parcel.
	Client int64 `json:"client"`
	// Status is the current status of the parcel.
	Status string `json:"status"`
	// Address is the destination address of the parcel.
	Address string `json:"address"`
	// CreatedAt is the timestamp of when the parcel was created.
	CreatedAt string `json:"created_at"`
}

// ParcelService provides operations for managing parcels.
//
// The ParcelService struct holds a reference to a ParcelStore,
// which is responsible for persisting and retrieving parcel data.
type ParcelService struct {
	// store is the interface for the underlying data storage
	// of parcels. It provides methods to create, read, update,
	// and delete parcel records.
	store ParcelStore
}

// NewParcelService creates a new instance of ParcelService.
//
// It takes a ParcelStore as a parameter, which is used to
// interface with the underlying data storage for parcel records.
// The function returns a ParcelService populated with the provided store.
func NewParcelService(store ParcelStore) ParcelService {
	return ParcelService{store: store}
}

// Register registers a new parcel with the given client ID and address.
//
// It creates a Parcel with the provided client ID and address,
// sets the status to ParcelStatusRegistered, and records the
// current time as the creation timestamp. The parcel is then
// added to the store, and its unique identifier is retrieved.
//
// If the addition to the store fails, an error is returned along
// with the partially created Parcel. If successful, the created
// Parcel, now with its assigned number, is returned along with
// a confirmation message logged to the standard output.
//
// Parameters:
//   - client: An integer representing the client ID associated
//     with the parcel.
//   - address: A string containing the destination address of
//     the parcel.
//
// Returns:
//   - The created Parcel, which includes the assigned number and
//     other details.
//   - An error, if any occurred during the registration process.
func (s ParcelService) Register(client int64, address string) (Parcel, error) {
	parcel := Parcel{
		Client:    client,
		Status:    ParcelStatusRegistered,
		Address:   address,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	err := s.store.Add(&parcel)
	if err != nil {
		return Parcel{}, err
	}

	fmt.Printf("Новая посылка № %d на адрес %s от клиента с идентификатором %d зарегистрирована %s\n",
		parcel.Number, parcel.Address, parcel.Client, parcel.CreatedAt)

	return parcel, nil
}

// PrintClientParcels prints the details of all parcels associated with a given client.
//
// It retrieves the parcels for the specified client by their ID using the
// ParcelStore's GetByClient method. If an error occurs during retrieval,
// it returns the error. Upon successfully fetching the parcels, it prints
// each parcel's details, including the parcel number, address, client ID,
// registration date, and status.
//
// Parameters:
// - client: An integer representing the client's unique identifier.
//
// Returns:
//   - An error, if any occurred during the retrieval process; otherwise,
//     it returns nil.
func (s ParcelService) PrintClientParcels(client int) error {
	parcels, err := s.store.GetByClient(client)
	if err != nil {
		return err
	}

	fmt.Printf("Посылки клиента %d:\n", client)
	for _, parcel := range parcels {
		fmt.Printf("Посылка № %d на адрес %s от клиента с идентификатором %d зарегистрирована %s, статус %s\n",
			parcel.Number, parcel.Address, parcel.Client, parcel.CreatedAt, parcel.Status)
	}

	return nil
}

// NextStatus updates the status of a parcel to its next logical state.
//
// It retrieves the parcel using the provided parcel number through the
// ParcelStore's Get method. If an error occurs during retrieval, it
// returns the error. Based on the current status of the parcel, it
// determines the next status in the sequence: from registered to sent,
// and from sent to delivered. If the parcel is already delivered,
// it simply returns nil without making any updates.
//
// If the status is successfully updated, it prints the parcel number
// and its new status. The new status is set using the ParcelStore's
// SetStatus method.
//
// Parameters:
// - number: An integer representing the unique identifier of the parcel.
//
// Returns:
//   - An error, if any occurred during retrieval or status update;
//     otherwise, it returns nil.
func (s ParcelService) NextStatus(number int) error {
	parcel, err := s.store.Get(number)
	if err != nil {
		return err
	}

	var nextStatus string
	switch parcel.Status {
	case ParcelStatusRegistered:
		nextStatus = ParcelStatusSent
	case ParcelStatusSent:
		nextStatus = ParcelStatusDelivered
	case ParcelStatusDelivered:
		return nil
	}

	fmt.Printf("У посылки № %d новый статус: %s\n", number, nextStatus)

	return s.store.SetStatus(number, nextStatus)
}

// ChangeAddress updates the delivery address of a parcel.
//
// This method changes the address of the parcel identified by its
// unique number. It calls the ParcelStore's SetAddress method to
// persist the new address in the storage system.
//
// Parameters:
//   - number: An integer representing the unique identifier of the parcel.
//   - address: A string containing the new address to which the parcel
//     should be sent.
//
// Returns:
// - An error if the address update fails; otherwise, it returns nil.
func (s ParcelService) ChangeAddress(number int, address string) error {
	return s.store.SetAddress(number, address)
}

// Delete removes a parcel from the store.
//
// This method deletes the parcel identified by its unique number from
// the storage system. It calls the ParcelStore's Delete method to
// perform the operation.
//
// Parameters:
// - number: An integer representing the unique identifier of the parcel.
//
// Returns:
// - An error if the deletion fails; otherwise, it returns nil.
func (s ParcelService) Delete(number int) error {
	return s.store.Delete(number)
}

// ParcelStore is a struct that represents the storage layer for parcels.
//
// It encapsulates a connection to the database, allowing for operations
// related to parcel records to be conducted. The `db` field is a pointer
// to an sql.DB instance, which is used to interact with the underlying
// database for adding, retrieving, and managing parcel data.
type ParcelStore struct {
	// db is a pointer to the SQL database connection.
	db *sql.DB
}

// NewParcelStore creates a new ParcelStore instance.
//
// This function initializes a new ParcelStore using the provided
// database connection. It returns a ParcelStore that can be used
// for operations on parcels.
//
// Parameters:
//   - db: A pointer to an sql.DB instance, representing the database
//     connection to be used by the ParcelStore.
//
// Returns:
// - A new instance of ParcelStore.
func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

// Add inserts a new parcel into the database and returns the newly created parcel's ID.
//
// Parameters:
// - p: the Parcel object containing the details of the parcel to be added.
//
// Returns:
// - The ID of the last inserted Parcel.
// - An error, if any occurs during the insert operation.
func (s ParcelStore) Add(p *Parcel) error {
	if p == nil {
		return errors.New("gotten pointer is equal to nil")
	}

	result, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)", p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return err
	}

	lastParcelID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	p.Number = lastParcelID

	return nil
}

// Get retrieves a parcel from the database by its number.
//
// Parameters:
// - number: the unique number of the parcel to retrieve.
//
// Returns:
// - The Parcel object corresponding to the given number.
// - An error, if any occurs during the retrieval operation.
func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE id = ?", number)

	gottenParcel := Parcel{}

	err := row.Scan(&gottenParcel.Number, &gottenParcel.Client, &gottenParcel.Status, &gottenParcel.Address, &gottenParcel.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Parcel{}, nil
	}

	if err != nil {
		return Parcel{}, err
	}

	return gottenParcel, nil
}

// GetByClient retrieves a list of parcels associated with a specific client.
//
// Parameters:
// - client: the unique identifier of the client whose parcels are to be retrieved.
//
// Returns:
// - A slice of Parcel objects corresponding to the given client.
// - An error, if any occurs during the retrieval operation.
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM percel WHERE client = ?", client)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var parcels []Parcel
	for rows.Next() {
		var newParcel Parcel

		err = rows.Scan(&newParcel.Number, &newParcel.Client, &newParcel.Status, &newParcel.Address, &newParcel.CreatedAt)
		if err != nil {
			return nil, err
		}

		parcels = append(parcels, newParcel)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return parcels, nil
}

// SetStatus updates the status of a parcel identified by its number.
//
// Parameters:
// - number: the unique number of the parcel to be updated.
// - status: the new status to set for the parcel.
//
// Returns:
// - An error, if any occurs during the update operation.
func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = ? WHERE number = ?", status, number)
	return err
}

// SetAddress updates the address of a parcel identified by its number.
//
// Parameters:
// - number: the unique number of the parcel to be updated.
// - address: the new address to set for the parcel.
//
// Returns:
// - An error, if any occurs during the update operation.
func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec("UPDATE parcel SET address = ? WHERE number = ?", address, number)
	return err
}

// Delete removes a parcel from the database identified by its number.
// The parcel will only be deleted if its status is 'registered'.
//
// Parameters:
// - number: the unique number of the parcel to be deleted.
//
// Returns:
// - An error, if any occurs during the deletion operation.
func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = ? AND status = registered", number)
	return err
}
