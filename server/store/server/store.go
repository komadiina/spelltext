package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/komadiina/spelltext/proto/store"
	"github.com/komadiina/spelltext/server/store/config"
	"github.com/komadiina/spelltext/utils/singleton/logging"
)

type StoreService struct {
	pb.UnimplementedStoreServer
	DbPool *pgxpool.Pool
	Config *config.Config
	Logger *logging.Logger
}

func tryConnect(s *StoreService, context context.Context, conninfo string, backoff time.Duration, maxRetries int, boFormula func(time.Duration) time.Duration) (pgx.Conn, error) {
	try := 1
	for {
		conn, err := pgx.Connect(context, conninfo)

		if err != nil && try >= maxRetries {
			// conn not established, max retries exceeded
			s.Logger.Fatal(err)
		} else if err == nil && try < maxRetries {
			// conn established within maxRetries
			s.Logger.Info("pgpool connection established")
			return *conn, nil
		} else if err != nil && try < maxRetries {
			// conn not established, backoff
			s.Logger.Warn("failed to establish database connection, backing off...", "reason", err, "backoff_seconds", backoff.Seconds())
			time.Sleep(backoff)
			backoff = boFormula(backoff)
			try++
		}
	}
}

func (s *StoreService) GetConn(ctx context.Context) *pgx.Conn {
	conninfo := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		s.Config.PgUser,
		s.Config.PgPass,
		s.Config.PgHost,
		s.Config.PgPort,
		s.Config.PgDbName,
		s.Config.PgSSLMode,
	)

	backoff := time.Second * 5 // secs
	time.Sleep(backoff)

	conn, err := tryConnect(s, ctx, conninfo, backoff, 5, func(backoff time.Duration) time.Duration {
		backoff = backoff + time.Second*5
		return backoff
	})

	if err != nil {
		return nil
	}

	return &conn
}

func (s *StoreService) ListVendors(ctx context.Context, req *pb.StoreListVendorRequest) (*pb.StoreListVendorResponse, error) {
	q := sq.Select("*").From("vendors")
	sql, _, err := q.ToSql()
	if err != nil {
		s.Logger.Error("failed to build query", "err", err)
		return nil, err
	}

	s.Logger.Info("running query", "query", sql)
	rows, err := s.DbPool.Query(ctx, sql)
	if err != nil {
		s.Logger.Error("failed to run query", "err", err)
		return nil, err
	}

	var vendors []*pb.Vendor
	for rows.Next() {
		v := &pb.Vendor{}
		err := rows.Scan(&v.VendorId, &v.VendorName, &v.VendorWareDescription)

		if err != nil {
			s.Logger.Error("failed to scan", "err", err)
			return nil, err
		}

		vendors = append(vendors, v)
	}

	return &pb.StoreListVendorResponse{Vendors: vendors}, nil
}

func (s *StoreService) ListVendorItems(ctx context.Context, req *pb.StoreListVendorItemRequest) (*pb.ListVendorItemResponse, error) {
	s.Logger.Infof("StoreListVendorItemRequest{%d}", req.GetVendorId())
	cte := sq.Select("v.id AS id, vw.item_type_id AS item_type_id").
		From("vendors AS v").
		InnerJoin("vendor_wares AS vw ON vw.vendor_id = v.id").
		Where("v.id = $1")

	cteSql, _, err := cte.ToSql()
	if err != nil {
		s.Logger.Error("failed to build cte", "err", err)
		return nil, err
	}
	prefix := fmt.Sprintf("WITH v_filt AS (%s)", cteSql)

	query := sq.
		Select("templ.*").
		From("item_templates AS templ").
		InnerJoin("v_filt ON v_filt.item_type_id = templ.item_type_id")

	sql, _, err := query.ToSql()
	if err != nil {
		s.Logger.Error("failed to build query", "err", err)
		return nil, err
	}

	sql = fmt.Sprintf("%s %s", prefix, sql)
	s.Logger.Info("running query", "query", strings.ReplaceAll(sql, "$1", fmt.Sprint(req.VendorId)))
	
	rows, err := s.DbPool.Query(ctx, sql, req.VendorId)
	if err != nil {
		s.Logger.Error("failed to query", "err", err)
		return nil, err
	}

	var items []*pb.Item
	for rows.Next() {
		it := &pb.Item{}
		err := rows.Scan(
			&it.Id,
			&it.Name,
			&it.ItemTypeId,
			&it.Rarity,
			&it.Stackable,
			&it.StackSize,
			&it.BindOnPickup,
			&it.Description,
			&it.Metadata,
		)
		if err != nil {
			s.Logger.Error("failed to scan", "err", err)
			return nil, err
		}
		items = append(items, it)
	}

	return &pb.ListVendorItemResponse{Items: items, TotalCount: -1}, nil
}

func (s *StoreService) AddItem(ctx context.Context, req *pb.AddItemRequest) (*pb.AddItemResponse, error) {
	panic("unimplemented")
}

func (s *StoreService) BuyItem(ctx context.Context, req *pb.BuyItemRequest) (*pb.BuyItemResponse, error) {
	panic("unimplemented")
}
