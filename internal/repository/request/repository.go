package request

import (
	"context"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	db "github.com/VadimGossip/platform_common/pkg/db/oracle"
	"github.com/sirupsen/logrus"

	"github.com/VadimGossip/drs_storage_tester/internal/model"
	def "github.com/VadimGossip/drs_storage_tester/internal/repository"
)

var _ def.RequestRepository = (*repository)(nil)

const (
	cdrTableName     = "cdr_based_drs_test_data"
	srcGwgrIdColumn  = "src_gwgr_id"
	anumberInColumn  = "anumber_in"
	bnumberInColumn  = "bnumber_in"
	anumberOutColumn = "anumber_out"
	bnumberOutColumn = "bnumber_out"
)

type repository struct {
	db db.Client
}

func NewRepository(db db.Client) *repository {
	return &repository{db: db}
}

func (r *repository) GetTaskRequests(ctx context.Context, limit int64) ([]model.TaskRequest, error) {
	cdrSelect := sq.Select(srcGwgrIdColumn,
		anumberInColumn,
		bnumberInColumn,
		anumberOutColumn,
		bnumberOutColumn).
		From(cdrTableName).
		PlaceholderFormat(sq.Colon)

	if limit > 0 {
		cdrSelect = cdrSelect.Where(sq.LtOrEq{"rownum": limit})
	}

	query, args, err := cdrSelect.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			logrus.WithFields(logrus.Fields{
				"handler": "CdrRequests",
				"problem": "rows close",
			}).Error(err)
		}
	}()

	var gwgrId int64
	var anumberIn, bnumberIn, anumberOut, bnumberOut uint64
	var anumberInStr, bnumberInStr, anumberOutStr, bnumberOutStr string
	result := make([]model.TaskRequest, 0)
	for rows.Next() {
		if err = rows.Scan(&gwgrId, &anumberInStr, &bnumberInStr, &anumberOutStr, &bnumberOutStr); err != nil {
			return nil, err
		}

		anumberIn, err = strconv.ParseUint(anumberInStr, 10, 64)
		if err != nil {
			continue
		}

		bnumberIn, err = strconv.ParseUint(bnumberInStr, 10, 64)
		if err != nil {
			continue
		}

		anumberOut, err = strconv.ParseUint(anumberOutStr, 10, 64)
		if err != nil {
			continue
		}

		bnumberOut, err = strconv.ParseUint(bnumberOutStr, 10, 64)
		if err != nil {
			continue
		}

		item := model.TaskRequest{
			GwgrId:      gwgrId,
			OrigAnumber: anumberIn,
			OrigBnumber: bnumberIn,
			Anumber:     anumberOut,
			Bnumber:     bnumberOut,
		}

		result = append(result, item)
	}

	return result, nil
}

func (r *repository) GetSupGwgrIds(ctx context.Context) ([]int64, error) {
	rows, err := r.db.DB().QueryContext(ctx, sqlSUPGWQuery)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			logrus.WithFields(logrus.Fields{
				"handler": "GetSupGwgrIds",
				"problem": "rows close",
			}).Error(err)
		}
	}()

	var gwgrId int64
	result := make([]int64, 0)
	for rows.Next() {
		if err = rows.Scan(&gwgrId); err != nil {
			return nil, err
		}
		result = append(result, gwgrId)
	}

	return result, nil
}
