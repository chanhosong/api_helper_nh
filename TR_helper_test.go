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

func TestF조회TR_NH(t *testing.T) {
	질의값 := new(lib.S질의값_단일종목)
	질의값.TR구분 = lib.TR조회
	질의값.TR코드 = lib.NH_TR_ETF_현재가_조회
	질의값.M종목코드 = lib.F임의_종목_ETF().G코드()

	lib.F체크포인트()

	응답_메시지 := F조회_NH(질의값)
	lib.F테스트_다름(t, 응답_메시지, nil)
	lib.F테스트_에러없음(t, 응답_메시지.G에러())
	lib.F테스트_같음(t, 응답_메시지.G길이(), 1)

	lib.F체크포인트()

	응답값 := lib.NewNH_ETF_현재가_조회_응답()
	lib.F테스트_에러없음(t, 응답_메시지.G값(0, 응답값))
	lib.F테스트_다름(t, 응답값, nil)

	lib.F체크포인트(응답값)
}

func TestTR실시간_서비스_등록_및_해지(t *testing.T) {
	if !lib.F한국증시_정규시장_거래시간임() {
		t.SkipNow()
	}

	lib.F테스트_참임(t, F접속됨_NH()) // NH서버에 접속 되었는 지 확인.

	// 실시간 정보 구독
	종목모음 := lib.F샘플_종목_모음_코스피200_ETF()
	종목코드_모음 := lib.F종목코드_추출(종목모음, len(종목모음))

	ch수신 := make(chan lib.I소켓_메시지, 10)
	lib.F테스트_에러없음(t, F실시간_정보_구독_NH(ch수신, lib.NH_RT코스피_체결, 종목코드_모음))

	// 실시간 정보 수신 확인
	for i := 0; i < 10; i++ {
		ch타임아웃 := time.After(lib.P1분)

		select {
		case 수신_메시지 := <-ch수신:
			lib.F테스트_에러없음(t, 수신_메시지.G에러())
		case <-ch타임아웃:
			lib.F문자열_출력("타임아웃. %v", i)
			t.FailNow()
		}
	}

	// 실시간 정보 해지
	lib.F테스트_에러없음(t, F실시간_정보_해지_NH(ch수신, lib.NH_RT코스피_체결, 종목코드_모음))
}
