package query

import (
	"database/sql"
	_ "fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL 驅動程式
)

// Query 函數動態處理 SQL 結果
func Query(db *sql.DB, SQLCode string) ([]map[string]interface{}, error) {
	log.Printf("Start querying from DB")

	rows, err := db.Query(SQLCode)
	if err != nil {
		log.Fatalf("Failed to query data from DB: %v", err)
		return nil, err
	}
	defer rows.Close()

	log.Println("Successfully queried from DB")

	// 取得欄位名稱
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 儲存結果的 slice
	var Results []map[string]interface{}

	for rows.Next() {
		// 建立一個 interface{} 切片來動態存儲數據
		Values := make([]interface{}, len(columns))
		ValuePointers := make([]interface{}, len(columns))

		// 建立指標陣列，讓 Scan 可以把值存進去
		for i := range Values {
			ValuePointers[i] = &Values[i]
		}

		// Scan 把數據填入 ValuePointers 指向的記憶體 , 意思等同於填入 Values 的儲存格中
		if err := rows.Scan(ValuePointers...); err != nil {
			return nil, err
		}

		// 轉成 map，key 是欄位名稱，value 是對應的數據
		RowMap := make(map[string]interface{})
		for i, col := range columns {
			// 需要處理 SQL 回傳的數據類型，確保 JSON 兼容
			RowMap[col] = Values[i]
		}

		Results = append(Results, RowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return Results, nil
}
