/*
 * Copyright (c) 2023 Yunshan Networks
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"context"
	"strconv"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/prompb"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/storage/remote"

	"github.com/deepflowio/deepflow/server/querier/app/prometheus/model"
)

type RemoteReadQuerierable struct {
	Args *model.PromQueryParams
	Ctx  context.Context
}

func (q *RemoteReadQuerierable) Querier(ctx context.Context, mint, maxt int64) (storage.Querier, error) {
	return &RemoteReadQuerier{Args: q.Args, Ctx: q.Ctx}, nil
}

type RemoteReadQuerier struct {
	Args *model.PromQueryParams
	Ctx  context.Context
}

// For PromQL instant query
func (q *RemoteReadQuerier) Select(sortSeries bool, hints *storage.SelectHints, matchers ...*labels.Matcher) storage.SeriesSet {
	startTimeS, err := (strconv.ParseFloat(q.Args.StartTime, 64))
	if err != nil {
		log.Error(err)
		return storage.ErrSeriesSet(err)
	}
	startTimeMs := int64(startTimeS * 1000)
	endTimeS, err := (strconv.ParseFloat(q.Args.EndTime, 64))
	if err != nil {
		log.Error(err)
		return storage.ErrSeriesSet(err)
	}
	endTimeMs := int64(endTimeS * 1000)
	prompbQuery, err := remote.ToQuery(startTimeMs, endTimeMs, matchers, hints)
	if err != nil {
		log.Error(err)
		return storage.ErrSeriesSet(err)
	}
	req := &prompb.ReadRequest{
		Queries:               []*prompb.Query{prompbQuery},
		AcceptedResponseTypes: []prompb.ReadRequest_ResponseType{prompb.ReadRequest_STREAMED_XOR_CHUNKS},
	}
	resp, err := promReaderExecute(req, q.Ctx)
	if err != nil {
		log.Error(err)
		return storage.ErrSeriesSet(err)
	}
	return remote.FromQueryResult(sortSeries, resp.Results[0])
}

func (q *RemoteReadQuerier) LabelValues(name string, matchers ...*labels.Matcher) ([]string, storage.Warnings, error) {
	return nil, nil, nil
}

func (q *RemoteReadQuerier) LabelNames(matchers ...*labels.Matcher) ([]string, storage.Warnings, error) {
	return nil, nil, nil
}

func (q *RemoteReadQuerier) Close() error {
	return nil
}

type RemoteReadRangeQuerierable struct {
	Args *model.PromQueryParams
	Ctx  context.Context
}

func (q *RemoteReadRangeQuerierable) Querier(ctx context.Context, mint, maxt int64) (storage.Querier, error) {
	return &RemoteReadRangeQuerier{Args: q.Args, Ctx: q.Ctx}, nil
}

type RemoteReadRangeQuerier struct {
	Args *model.PromQueryParams
	Ctx  context.Context
}

// For PromQL range query
func (q *RemoteReadRangeQuerier) Select(sortSeries bool, hints *storage.SelectHints, matchers ...*labels.Matcher) storage.SeriesSet {
	startS, err := (strconv.ParseFloat(q.Args.StartTime, 64))
	if err != nil {
		log.Error(err)
		return storage.ErrSeriesSet(err)
	}
	endS, err := (strconv.ParseFloat(q.Args.EndTime, 64))
	if err != nil {
		log.Error(err)
		return storage.ErrSeriesSet(err)
	}
	startMs := int64(startS * 1000)
	endMs := int64(endS * 1000)
	prompbQuery, err := remote.ToQuery(startMs, endMs, matchers, hints)
	if err != nil {
		log.Error(err)
		return storage.ErrSeriesSet(err)
	}
	req := &prompb.ReadRequest{
		Queries:               []*prompb.Query{prompbQuery},
		AcceptedResponseTypes: []prompb.ReadRequest_ResponseType{prompb.ReadRequest_STREAMED_XOR_CHUNKS},
	}
	resp, err := promReaderExecute(req, q.Ctx)
	if err != nil {
		log.Error(err)
		return storage.ErrSeriesSet(err)
	}
	return remote.FromQueryResult(sortSeries, resp.Results[0])
}

func (q *RemoteReadRangeQuerier) LabelValues(name string, matchers ...*labels.Matcher) ([]string, storage.Warnings, error) {
	return nil, nil, nil
}

func (q *RemoteReadRangeQuerier) LabelNames(matchers ...*labels.Matcher) ([]string, storage.Warnings, error) {
	return nil, nil, nil
}

func (q *RemoteReadRangeQuerier) Close() error {
	return nil
}
