package main

import (
	"context"
	"github.com/HuXin0817/bitkv"
	"github.com/HuXin0817/bitkv/interval/errors"
	"github.com/HuXin0817/bitkv/interval/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetLogic struct {
	ctx    context.Context
	svcCtx *ServiceContext
	logx.Logger
}

func NewGetLogic(ctx context.Context, svcCtx *ServiceContext) *GetLogic {
	return &GetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLogic) Get(in *pb.Key) (*pb.Value, error) {
	if err := errors.CheckKey(in.Key); err != nil {
		return nil, err
	}
	return &pb.Value{Value: bitkv.Get(in.Key)}, nil
}
