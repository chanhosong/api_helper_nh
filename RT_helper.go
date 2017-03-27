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
	"time"
	"github.com/go-mangos/mangos"
	"github.com/boltdb/bolt"
)

var 소켓SUB_NH실시간_정보 mangos.Socket
var 실시간_정보_중계_Go루틴_실행_중 = lib.New안전한_bool(false)
var nh구독채널_저장소 = lib.New구독채널_저장소()

func go루틴_실시간_정보_중계(ch초기화 chan lib.T신호) {
	lib.F메모("모든 정보를 배포하는 대신, 채널별로 구독한 정보만 배포하는 방안을 생각해 볼 것.")

	if 에러 := 실시간_정보_중계_Go루틴_실행_중.S값(true); 에러 != nil {
		ch초기화 <- lib.P신호_초기화
		return
	}

	defer 실시간_정보_중계_Go루틴_실행_중.S값(false)

	nh대기_중_데이터_저장소 := new대기_중_데이터_저장소()
	ch종료 := lib.F공통_종료_채널()
	ch초기화 <- lib.P신호_초기화

	for {
		바이트_모음, 에러 := 소켓SUB_NH실시간_정보.Recv()
		if 에러 != nil {
			switch 에러.Error() {
			case "connection closed":
				소켓SUB_NH실시간_정보, 에러 = lib.New소켓SUB(lib.P주소_NH_실시간_CBOR)
			}

			lib.F에러_출력(에러)

			continue
		}

		소켓_메시지 := lib.New소켓_메시지from바이트_모음(바이트_모음)
		if 소켓_메시지.G에러() != nil {
			lib.F에러_출력(소켓_메시지.G에러())
			continue
		}

		for _, ch수신 := range nh구독채널_저장소.G채널_모음() {
			f실시간_정보_중계_도우미(ch수신, 소켓_메시지, nh구독채널_저장소, nh대기_중_데이터_저장소)
		}

		nh대기_중_데이터_저장소.S재전송()

		// 종료 조건 확인
		select {
		case <-ch종료:
			return
		default:    // 종료 신호 없으면 반복할 것.
		}

		lib.F실행권한_양보()
	}
}

// 이미 닫혀 있는 채널에 전송하면 패닉이 발생하는 데,
// 이 패닉을 적절하게 처리하도록 하기 위해서 해당 부분을 별도의 함수로 분리함.
func f실시간_정보_중계_도우미(ch수신 chan lib.I소켓_메시지, 소켓_메시지 lib.I소켓_메시지,
	nh구독채널_저장소 *lib.S구독채널_저장소, nh대기_중_데이터_저장소 *s대기_중_데이터_저장소) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{
		M함수with패닉내역: func(r interface{}) {
			lib.F에러_출력(r)
			// 이미 닫혀있는 채널에 전송하면 패닉 발생하므로, 해당 채널을 구독 채널 목록에서 제외시킴.
			nh구독채널_저장소.S삭제(ch수신)
		}})

	select {
	case ch수신 <- 소켓_메시지: // 메시지 중계 성공
	default: // 메시지 중계 실패.
		nh대기_중_데이터_저장소.S추가(소켓_메시지, ch수신)    // 저장소에 보관하고 추후 재전송
	}
}

var nh_ETF_실시간_데이터_수집 = lib.New안전한_bool(false)
var 버킷ID_호가_잔량 = []byte(P버킷ID_NH호가_잔량)
var 버킷ID_시간외_호가_잔량 = []byte(P버킷ID_NH시간외_호가잔량)
var 버킷ID_예상_호가_잔량 = []byte(P버킷ID_NH예상_호가잔량)
var 버킷ID_체결 = []byte(P버킷ID_NH체결)
var 버킷ID_ETF_NAV = []byte(P버킷ID_NH_ETF_NAV)
var 버킷ID_업종지수 = []byte(P버킷ID_NH업종지수)

func go루틴_실시간_데이터_저장(ch초기화 chan lib.T신호, ch수신 chan lib.I소켓_메시지, db lib.I데이터베이스) {
	if 에러 := nh_ETF_실시간_데이터_수집.S값(true); 에러 != nil {
		ch초기화 <- lib.P신호_초기화
		return
	}

	defer nh_ETF_실시간_데이터_수집.S값(false)

	var 수신_메시지 lib.I소켓_메시지
	ch종료 := lib.F공통_종료_채널()

	ch초기화 <- lib.P신호_초기화

	// 실시간 데이터 수신
	for {
		select {
		case 수신_메시지 = <-ch수신:
			//lib.F체크포인트("데이터 수신")
			if 수신_메시지.G에러() != nil {
				lib.F에러_출력(수신_메시지.G에러())
				continue
			}
		case <-ch종료:
			return
		}

		// DB에 수신값 저장
		if 에러 := fNH_실시간_데이터_저장_도우미(수신_메시지, db); 에러 != nil {
			if 에러.Error() == bolt.ErrDatabaseNotOpen.Error() &&
				lib.F테스트_모드_실행_중() {
				// 테스트 과정에서 DB를 닫은 것이니 종료.
				return
			}

			lib.F체크포인트("데이터 저장 에러")
			lib.F문자열_출력("에러 발생~!\n%s\n~!에러 끝", 에러.Error())
		}

		//lib.F체크포인트("데이터 저장 성공")
	}
}

func fNH_실시간_데이터_저장_도우미(수신_메시지 lib.I소켓_메시지, 데이터베이스 lib.I데이터베이스) (에러 error) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M에러: &에러})

	lib.F조건부_패닉(데이터베이스 == nil, "nil 데이터베이스")
	lib.F조건부_패닉(수신_메시지 == nil, "nil 메시지")
	lib.F조건부_패닉(수신_메시지.G길이() != 1, "예상하지 못한 메시지 길이. %v", 수신_메시지.G길이())

	질의 := new(lib.S데이터베이스_질의)

	switch 수신_메시지.G자료형_문자열(0) {
	case P버킷ID_NH호가_잔량:
		s := new(lib.NH호가_잔량)
		lib.F에러2패닉(수신_메시지.G값(0, s))

		시각, 에러 := s.M시각.MarshalBinary()
		lib.F에러2패닉(에러)

		질의.M버킷ID = 버킷ID_호가_잔량
		질의.M키 = append([]byte(s.M종목코드), 시각...)
	case P버킷ID_NH시간외_호가잔량:
		s := new(lib.NH시간외_호가잔량)
		lib.F에러2패닉(수신_메시지.G값(0, s))

		시각, 에러 := s.M시각.MarshalBinary()
		lib.F에러2패닉(에러)

		질의.M버킷ID = 버킷ID_시간외_호가_잔량
		질의.M키 = append([]byte(s.M종목코드), 시각...)
	case P버킷ID_NH예상_호가잔량:
		s := new(lib.NH예상_호가잔량)
		lib.F에러2패닉(수신_메시지.G값(0, s))

		시각, 에러 := s.M시각.MarshalBinary()
		lib.F에러2패닉(에러)

		질의.M버킷ID = 버킷ID_예상_호가_잔량
		질의.M키 = append([]byte(s.M종목코드), 시각...)
	case P버킷ID_NH체결:
		s := new(lib.NH체결)
		lib.F에러2패닉(수신_메시지.G값(0, s))

		시각, 에러 := s.M시각.MarshalBinary()
		lib.F에러2패닉(에러)

		질의.M버킷ID = 버킷ID_체결
		질의.M키 = append([]byte(s.M종목코드), 시각...)
	case P버킷ID_NH_ETF_NAV:
		s := new(lib.NH_ETF_NAV)
		lib.F에러2패닉(수신_메시지.G값(0, s))

		시각, 에러 := s.M시각.MarshalBinary()
		lib.F에러2패닉(에러)

		질의.M버킷ID = 버킷ID_ETF_NAV
		질의.M키 = append([]byte(s.M종목코드), 시각...)
	case P버킷ID_NH업종지수:
		s := new(lib.NH업종지수)
		lib.F에러2패닉(수신_메시지.G값(0, s))

		시각, 에러 := s.M시각.MarshalBinary()
		lib.F에러2패닉(에러)

		질의.M버킷ID = 버킷ID_업종지수
		질의.M키 = append([]byte(s.M업종코드), 시각...)
	default:
		lib.F패닉("예상하지 못한 자료형. %v", 수신_메시지.G자료형_문자열(0))
	}

	질의.M값, 에러 = 수신_메시지.G바이트_모음(0)
	lib.F에러2패닉(에러)

	return 데이터베이스.S업데이트(질의)
}

func F실시간_데이터_수집_NH_ETF(파일명 string, 종목코드_모음 []string) (db lib.I데이터베이스, 에러 error) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{
		M에러: &에러,
		M함수with패닉내역: func(r interface{}) {
			lib.New에러with출력(r)
			db = nil
		}})

	ch수신 := make(chan lib.I소켓_메시지, 10000)
	F실시간_데이터_구독_NH_ETF(ch수신, 종목코드_모음)

	ch초기화 := make(chan lib.T신호)
	go go루틴_실시간_정보_중계(ch초기화)
	신호 := <-ch초기화
	lib.F조건부_패닉(신호 != lib.P신호_초기화, "예상하지 못한 신호. %v", 신호)

	ch결과물 := make(chan interface{}, 1)
	ch타임아웃 := time.After(lib.P10초)
	go f실시간_데이터_수집_NH_ETF_도우미(파일명, ch결과물)

	select {
	case 결과물 := <-ch결과물:
		switch 결과값 := 결과물.(type) {
		case error:
			lib.F패닉(결과값)
			return nil, 결과값
		case lib.I데이터베이스:
			db = 결과값
			break
		default:
			lib.F패닉("예상하지 못한 자료형")
			return nil, nil
		}

		break
	case <-ch타임아웃:
		lib.F패닉("타임아웃")
	}

	go go루틴_실시간_데이터_저장(ch초기화, ch수신, db)

	신호 = <-ch초기화

	lib.F조건부_패닉(신호 != lib.P신호_초기화, "예상하지 못한 신호. %v", 신호)

	//lib.F체크포인트("데이터 저장 초기화 완료.")

	// go루틴 종료는 'lib.F공통_종료_채널_닫은_후_재설정()'으로 한다.

	return db, nil
}

func f실시간_데이터_수집_NH_ETF_도우미(파일명 string, ch결과물 chan interface{}) {
	db, 에러 := lib.NewBoltDB(파일명)

	switch {
	case 에러 != nil:
		ch결과물 <- 에러
	default:
		ch결과물 <- db
	}
}

func F실시간_데이터_구독_NH_ETF(ch수신 chan lib.I소켓_메시지, 종목코드_모음 []string) (에러 error) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M에러: &에러})

	lib.F에러2패닉(F실시간_정보_구독_NH(ch수신, lib.NH_RT코스피_호가_잔량, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	lib.F에러2패닉(F실시간_정보_구독_NH(ch수신, lib.NH_RT코스피_체결, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	lib.F에러2패닉(F실시간_정보_구독_NH(ch수신, lib.NH_RT코스피_ETF_NAV, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	lib.F에러2패닉(F실시간_정보_구독_NH(ch수신, lib.NH_RT코스피_시간외_호가_잔량, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	lib.F에러2패닉(F실시간_정보_구독_NH(ch수신, lib.NH_RT코스피_예상_호가_잔량, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	return nil
}

func F실시간_데이터_해지_NH_ETF(종목코드_모음 []string) (에러 error) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M에러: &에러})

	lib.F에러2패닉(F실시간_정보_해지_NH(lib.NH_RT코스피_호가_잔량, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	lib.F에러2패닉(F실시간_정보_해지_NH(lib.NH_RT코스피_체결, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	lib.F에러2패닉(F실시간_정보_해지_NH(lib.NH_RT코스피_ETF_NAV, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	lib.F에러2패닉(F실시간_정보_해지_NH(lib.NH_RT코스피_시간외_호가_잔량, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	lib.F에러2패닉(F실시간_정보_해지_NH(lib.NH_RT코스피_예상_호가_잔량, 종목코드_모음))
	lib.F대기(lib.P500밀리초)

	return nil
}