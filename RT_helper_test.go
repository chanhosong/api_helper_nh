/* Copyright (C) 2015-2016 김운하(UnHa Kim)  unha.kim@kuh.pe.kr

이 파일은 GHTS의 일부입니다.

이 프로그램은 자유 소프트웨어입니다.
소프트웨어의 피양도자는 자유 소프트웨어 재단이 공표한 GNU LGPL 2.1판
규정에 따라 프로그램을 개작하거나 재배포할 수 있습니다.

이 프로그램은 유용하게 사용될 수 있으리라는 희망에서 배포되고 있지만,
특정한 목적에 적합하다거나, 이익을 안겨줄 수 있다는 묵시적인 보증을 포함한
어떠한 형태의 보증도 제공하지 않습니다.
보다 자세한 사항에 대해서는 GNU LGPL 2.1판을 참고하시기 바랍니다.
GNU LGPL 2.1판은 이 프로그램과 함께 제공됩니다.
만약, 이 문서가 누락되어 있다면 자유 소프트웨어 재단으로 문의하시기 바랍니다.
(자유 소프트웨어 재단 : Free Software Foundation, Inc.,
59 Temple Place - Suite 330, Boston, MA 02111-1307, USA)

Copyright (C) 2015년 UnHa Kim (unha.kim@kuh.pe.kr)

This file is part of GHTS.

GHTS is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, version 2.1 of the License.

GHTS is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with GHTS.  If not, see <http://www.gnu.org/licenses/>. */

package api_helper_nh

import (
	"github.com/ghts/lib"
	"testing"
	"time"
)

// 여기에서 별도로 도우미 함수 기능을 테스트 하지 않고,
// 도우미 함수의 테스트는 메인 함수의 기능 테스트로 대체.

func TestF자료형_문자열(t *testing.T) {
	lib.F테스트_같음(t, P버킷ID_NH호가_잔량, lib.F자료형_문자열(lib.NH호가_잔량{}))
	lib.F테스트_같음(t, P버킷ID_NH시간외_호가잔량, lib.F자료형_문자열(lib.NH시간외_호가잔량{}))
	lib.F테스트_같음(t, P버킷ID_NH예상_호가잔량, lib.F자료형_문자열(lib.NH예상_호가잔량{}))
	lib.F테스트_같음(t, P버킷ID_NH체결, lib.F자료형_문자열(lib.NH체결{}))
	lib.F테스트_같음(t, P버킷ID_NH_ETF_NAV, lib.F자료형_문자열(lib.NH_ETF_NAV{}))
	lib.F테스트_같음(t, P버킷ID_NH업종지수, lib.F자료형_문자열(lib.NH업종지수{}))
}

func TestGo루틴_실시간_정보_중계(t *testing.T) {
	ch실시간_데이터 := make(chan lib.I소켓_메시지, 100)

	종목_모음, 에러 := lib.F종목모음_코스피()
	lib.F테스트_에러없음(t, 에러)

	종목코드_모음 := lib.F종목코드_추출(종목_모음, 100)

	lib.F테스트_에러없음(t, F실시간_정보_구독_NH(ch실시간_데이터, lib.NH_RT코스피_호가_잔량, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	lib.F테스트_에러없음(t, F실시간_정보_구독_NH(ch실시간_데이터, lib.NH_RT코스피_체결, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	lib.F테스트_에러없음(t, F실시간_정보_구독_NH(ch실시간_데이터, lib.NH_RT코스피_ETF_NAV, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	lib.F테스트_에러없음(t, F실시간_정보_구독_NH(ch실시간_데이터, lib.NH_RT코스피_시간외_호가_잔량, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	lib.F테스트_에러없음(t, F실시간_정보_구독_NH(ch실시간_데이터, lib.NH_RT코스피_예상_호가_잔량, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	ch초기화 := make(chan lib.T신호)
	go go루틴_실시간_정보_중계(ch초기화)
	lib.F테스트_같음(t, <-ch초기화, lib.P신호_초기화)

	ch대기시간_초과 := time.After(lib.P30초)

	var 테스트_수량 int
	if lib.F한국증시_정규시장_거래시간임() {
		테스트_수량 = 10
	} else {
		테스트_수량 = 1
	}

	수신_수량 := 0

	for {
		select {
		case 소켓_메시지 := <-ch실시간_데이터:
			lib.F테스트_같음(t, 소켓_메시지.G자료형_문자열(0),
				P버킷ID_NH호가_잔량, P버킷ID_NH시간외_호가잔량, P버킷ID_NH예상_호가잔량, P버킷ID_NH체결, P버킷ID_NH_ETF_NAV, P버킷ID_NH업종지수)

			if 수신_수량++; 수신_수량 > 테스트_수량 {
				return
			}
		case <-ch대기시간_초과:
			lib.F체크포인트("대기시간 초과")
			t.FailNow()
		}
	}
}

func TestGo루틴_실시간_데이터_저장(t *testing.T) {
	var db lib.I데이터베이스
	파일명 := f테스트용_실시간_데이터_파일명()

	defer func() {
		if db != nil {
			db.S종료()
		}}()

	ch실시간_데이터 := make(chan lib.I소켓_메시지, 100)

	종목_모음, 에러 := lib.F종목모음_코스피()
	lib.F테스트_에러없음(t, 에러)

	종목코드_모음 := lib.F종목코드_추출(종목_모음, 20)

	lib.F테스트_에러없음(t, F실시간_정보_구독_NH(ch실시간_데이터, lib.NH_RT코스피_호가_잔량, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	lib.F테스트_에러없음(t, F실시간_정보_구독_NH(ch실시간_데이터, lib.NH_RT코스피_체결, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	lib.F테스트_에러없음(t, F실시간_정보_구독_NH(ch실시간_데이터, lib.NH_RT코스피_ETF_NAV, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	lib.F테스트_에러없음(t, F실시간_정보_구독_NH(ch실시간_데이터, lib.NH_RT코스피_시간외_호가_잔량, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	lib.F테스트_에러없음(t, F실시간_정보_구독_NH(ch실시간_데이터, lib.NH_RT코스피_예상_호가_잔량, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	db, 에러 = lib.NewBoltDB(파일명)
	lib.F에러2패닉(에러)

	ch초기화 := make(chan lib.T신호)

	go go루틴_실시간_정보_중계(ch초기화)
	lib.F테스트_같음(t, <-ch초기화, lib.P신호_초기화)

	go go루틴_실시간_데이터_저장(ch초기화, ch실시간_데이터, db)
	lib.F테스트_같음(t, <-ch초기화, lib.P신호_초기화)

	ch대기시간_초과 := time.After(lib.P10초)

	var 테스트_수량 int
	if lib.F한국증시_정규시장_거래시간임() {
		테스트_수량 = 10
	} else {
		테스트_수량 = 1
	}

	for {
		lib.F대기(lib.P1초)

		저장_수량 := db.G수량in버킷(버킷ID_호가_잔량) +
			db.G수량in버킷(버킷ID_시간외_호가_잔량) +
			db.G수량in버킷(버킷ID_예상_호가_잔량) +
			db.G수량in버킷(버킷ID_체결) +
			db.G수량in버킷(버킷ID_ETF_NAV) +
			db.G수량in버킷(버킷ID_업종지수)

		if 저장_수량 >= 테스트_수량 {
			return
		}

		select {
		case <-ch대기시간_초과:
			t.FailNow()
		default:    // PASS
		}
	}
}

func TestF실시간_데이터_수집_NH_ETF(t *testing.T) {
	var db lib.I데이터베이스
	파일명 := f테스트용_실시간_데이터_파일명()

	defer func() {
		if db != nil {
			db.S종료()
		}}()

	종목_모음 := []*lib.S종목{
		lib.New종목("069500", "KODEX 200", lib.P시장구분_ETF),
		lib.New종목("114800", "KODEX 인버스", lib.P시장구분_ETF),
		lib.New종목("122630", "KODEX 레버리지", lib.P시장구분_ETF),
		lib.New종목("252670", "KODEX 200 선물인버스2X", lib.P시장구분_ETF),
		lib.New종목("069660", "KOSEF 200", lib.P시장구분_ETF),
		lib.New종목("152280", "KOSEF 200선물", lib.P시장구분_ETF),
		lib.New종목("253250", "KOSEF 200 선물레버리지", lib.P시장구분_ETF),
		lib.New종목("253240", "KOSEF 200 선물인버스", lib.P시장구분_ETF),
		lib.New종목("253230", "KOSEF 200 선물인버스2X", lib.P시장구분_ETF),
		lib.New종목("102110", "TIGER 200", lib.P시장구분_ETF),
		lib.New종목("252710", "TIGER 200 선물인버스2X", lib.P시장구분_ETF),
		lib.New종목("105190", "KINDEX 200", lib.P시장구분_ETF),
		lib.New종목("108590", "TREX 200", lib.P시장구분_ETF),
		lib.New종목("148020", "KBSTAR 200", lib.P시장구분_ETF),
		lib.New종목("252400", "KBSTAR 200 선물레버리지", lib.P시장구분_ETF),
		lib.New종목("252410", "KBSTAR 200 선물인버스", lib.P시장구분_ETF),
		lib.New종목("252420", "KBSTAR 200 선물인버스2X", lib.P시장구분_ETF),
		lib.New종목("152100", "ARIRANG 200", lib.P시장구분_ETF),
		lib.New종목("253150", "ARIRANG 200 선물레버리지", lib.P시장구분_ETF),
		lib.New종목("253160", "ARIRANG 200 선물인버스2X", lib.P시장구분_ETF)}

	종목코드_모음 := lib.F종목코드_추출(종목_모음, 20)

	db, 에러 := F실시간_데이터_수집_NH_ETF(파일명, 종목코드_모음)
	lib.F테스트_에러없음(t, 에러)
	lib.F테스트_다름(t, db, nil)

	var 테스트_수량 int
	if lib.F한국증시_정규시장_거래시간임() {
		테스트_수량 = 10
	} else {
		테스트_수량 = 1
	}

	제한시간 := time.Now().Add(lib.P3분)

	for {
		lib.F테스트_다름(t, db, nil)

		저장_수량 := db.G수량in버킷(버킷ID_호가_잔량) +
			db.G수량in버킷(버킷ID_시간외_호가_잔량) +
			db.G수량in버킷(버킷ID_예상_호가_잔량) +
			db.G수량in버킷(버킷ID_체결) +
			db.G수량in버킷(버킷ID_ETF_NAV) +
			db.G수량in버킷(버킷ID_업종지수)

		switch {
		case 저장_수량 >= 테스트_수량:
				lib.F문자열_출력("OK")
				return
		case time.Now().After(제한시간):
				t.FailNow()
		default:
			lib.F대기(lib.P1초)
		}
	}
}

func f테스트용_실시간_데이터_파일명() string {
	if lib.F테스트_모드_실행_중() {
		// 반복된 테스트로 인한 파일명 중복을 피하기 위함.
		return "test_realtime_data_NH_" + time.Now().Format("20060102_150405") + ".dat"
	}

	return "realtime_data_NH_" + time.Now().Format(lib.P일자_형식) + ".dat"
}