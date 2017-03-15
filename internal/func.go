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

package internal

import (
	"github.com/ghts/lib"
	ps "github.com/mitchellh/go-ps"
	"strings"
	"sync"
	"time"
)

var _NH_API커넥터_경로 = lib.F_GOPATH() + `/src/github.com/ghts/api_bridge_nh/api_bridge_nh.exe`
var _NH_API커넥터_실행_잠금 = new(sync.Mutex)

func F_NH_API커넥터_실행() {
	_NH_API커넥터_실행_잠금.Lock()
	defer _NH_API커넥터_실행_잠금.Unlock()

	if fNH_API커넥터_실행_중() {
		return
	}

	_, 에러 := lib.F외부_프로세스_실행(_NH_API커넥터_경로)
	lib.F에러2패닉(에러)
}

func fNH_API커넥터_실행_중() (실행_중 bool) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수: func() { 실행_중 = false }})

	프로세스_모음, 에러 := ps.Processes()
	lib.F에러2패닉(에러)

	for _, 프로세스 := range 프로세스_모음 {
		if 실행화일명 := 프로세스.Executable(); strings.HasSuffix(_NH_API커넥터_경로, 실행화일명) {
			return true
		}
	}

	return false
}

func F접속_NH() (로그인_정보 *lib.NH로그인_정보, 에러 error) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{
		M에러: &에러,
		M함수: func() { 로그인_정보 = nil }})

	질의값 := new(lib.S질의값_단순TR)
	질의값.TR구분 = lib.TR접속

	응답 := lib.New소켓_질의(lib.P주소_NH_TR, lib.CBOR, lib.P30초).S질의(질의값).G응답()
	lib.F에러2패닉(응답.G에러())
	lib.F에러2패닉(응답.G값(0, &로그인_정보))

	return
}

// 접속 되었는 지 확인.
func F접속됨_NH() (참거짓 bool) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수: func() { 참거짓 = false }})

	질의값 := new(lib.S질의값_단순TR)
	질의값.TR구분 = lib.TR접속됨

	응답 := lib.New소켓_질의(lib.P주소_NH_TR, lib.CBOR, lib.P30초).S질의(질의값).G응답()
	lib.F에러2패닉(응답.G에러())
	lib.F에러2패닉(응답.G값(0, &참거짓))

	return 참거짓
}

func fNH_실시간_데이터_파일명() string {
	if lib.F테스트_모드_실행_중() {
		// 반복된 테스트로 인한 파일명 중복을 피하기 위함.
		return "test_realtime_data_NH_" + time.Now().Format("20060102_150405") + ".dat"
	}

	return "realtime_data_NH_" + time.Now().Format(lib.P일자_형식) + ".dat"
}