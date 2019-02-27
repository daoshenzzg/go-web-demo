package demo

import (
	"database/sql"

	"person.mgtv.com/framework/database"
	demoModel "person.mgtv.com/model/demo"
)

type DemoDao struct {
	DB *sql.DB
}

func NewDemoDao() *DemoDao {
	return &DemoDao{
		DB: database.Get("db.feed"),
	}
}

func (dao *DemoDao) GetFeed(feedId string) (feed *demoModel.MaxTimeline, err error) {
	query := `SELECT id, uid, title FROM max_timeline WHERE id = ?`
	rows, err := dao.DB.Query(query, feedId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		feed = new(demoModel.MaxTimeline)
		err = rows.Scan(&feed.Id, &feed.Uid, &feed.Title)
	}

	return feed, err
}

func (dao *DemoDao) UpdateSeq(tx *sql.Tx, id int64, seq string) (err error) {
	query := `UPDATE max_sync_seq SET seq = ? WHERE id = ?`
	_, err = tx.Exec(query, seq, id)

	return err
}
