package main

import (
	"context"

	"github.com/HuXin0817/bitkv"
	"github.com/HuXin0817/bitkv/interval/errors"
	"github.com/HuXin0817/bitkv/interval/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type PutLogic struct {
	ctx    context.Context
	svcCtx *ServiceContext
	logx.Logger
}

func NewPutLogic(ctx context.Context, svcCtx *ServiceContext) *PutLogic {
	return &PutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PutLogic) Put(in *pb.Pair) (*pb.Nil, error) {
	if err := errors.CheckKey(in.Key); err != nil {
		return nil, err
	}
	if err := errors.CheckValue(in.Value); err != nil {
		return nil, err
	}
	bitkv.Put(in.Key, in.Value)
	return &pb.Nil{}, nil
}
