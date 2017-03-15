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

package main

import (
	"github.com/ghts/lib"
	nh "github.com/ghts/api_helper_nh"
	"time"
)

func main() {
	var 에러 error

	defer lib.F에러패닉_처리(lib.S에러패닉_처리{
		M에러: &에러,
		M함수 : func() { lib.F에러_출력(에러) }})

	lib.F에러2패닉(nh.F접속_NH())

	종목_모음 := lib.F샘플_종목_모음_코스피200_ETF()
	종목코드_모음 := lib.F종목코드_추출(종목_모음, len(종목_모음))
	파일명 := "RealTimeData_NH_" + time.Now().Format(lib.P일자_형식) + ".dat"

	db, 에러 := nh.F실시간_데이터_수집_NH_ETF(파일명, 종목코드_모음)
	lib.F에러2패닉(에러)

	// 1분마다 저장 수량 출력
	버킷ID_호가_잔량 := []byte(nh.P버킷ID_NH호가_잔량)
	버킷ID_시간외_호가_잔량 := []byte(nh.P버킷ID_NH시간외_호가잔량)
	버킷ID_예상_호가_잔량 := []byte(nh.P버킷ID_NH예상_호가잔량)
	버킷ID_체결 := []byte(nh.P버킷ID_NH체결)
	버킷ID_ETF_NAV := []byte(nh.P버킷ID_NH_ETF_NAV)
	버킷ID_업종지수 := []byte(nh.P버킷ID_NH업종지수)

	for {
		대기 := time.After(lib.P1분)
		select {
		case <-대기:
		}

		저장_수량 := db.G수량in버킷(버킷ID_호가_잔량) +
			db.G수량in버킷(버킷ID_시간외_호가_잔량) +
			db.G수량in버킷(버킷ID_예상_호가_잔량) +
			db.G수량in버킷(버킷ID_체결) +
			db.G수량in버킷(버킷ID_ETF_NAV) +
			db.G수량in버킷(버킷ID_업종지수)

		lib.F문자열_출력("%s : %v", time.Now().Format(lib.P간략한_시간_형식), 저장_수량)
	}
}
