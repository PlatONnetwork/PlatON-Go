(module
 (type $FUNCSIG$vii (func (param i32 i32)))
 (type $FUNCSIG$vi (func (param i32)))
 (type $FUNCSIG$vj (func (param i64)))
 (type $FUNCSIG$v (func))
 (import "env" "printi" (func $printi (param i64)))
 (import "env" "prints" (func $prints (param i32)))
 (import "env" "prints_l" (func $prints_l (param i32 i32)))
 (import "env" "__cxa_pure_virtual" (func $__cxa_pure_virtual))
 (table 3 3 anyfunc)
 (elem (i32.const 0) $__wasm_nullptr $_ZN6platon5Token4initEv $__importThunk___cxa_pure_virtual)
 (memory $0 2)
 (data (i32.const 4) "P\90\01\00")
 (data (i32.const 16) "from:% to:% asset: % \n\00")
 (data (i32.const 48) "sdf\00")
 (data (i32.const 52) "\00\00\00\00\00\00\00\00\01\00\00\00")
 (data (i32.const 64) "\00\00\00\00\00\00\00\00\02\00\00\00")
 (export "memory" (memory $0))
 (export "_ZeqRK11checksum256S1_" (func $_ZeqRK11checksum256S1_))
 (export "_ZeqRK11checksum160S1_" (func $_ZeqRK11checksum160S1_))
 (export "_ZneRK11checksum160S1_" (func $_ZneRK11checksum160S1_))
 (export "transfer" (func $transfer))
 (export "atoi" (func $atoi))
 (export "memcmp" (func $memcmp))
 (export "strlen" (func $strlen))
 (func $_ZeqRK11checksum256S1_ (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (set_local $2
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 16)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $2)
  )
  (i32.store offset=12
   (get_local $2)
   (get_local $0)
  )
  (i32.store offset=8
   (get_local $2)
   (get_local $1)
  )
  (set_local $1
   (call $memcmp
    (i32.load offset=12
     (get_local $2)
    )
    (i32.load offset=8
     (get_local $2)
    )
    (i32.const 32)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $2)
    (i32.const 16)
   )
  )
  (i32.eq
   (get_local $1)
   (i32.const 0)
  )
 )
 (func $_ZeqRK11checksum160S1_ (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (set_local $2
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 16)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $2)
  )
  (i32.store offset=12
   (get_local $2)
   (get_local $0)
  )
  (i32.store offset=8
   (get_local $2)
   (get_local $1)
  )
  (set_local $1
   (call $memcmp
    (i32.load offset=12
     (get_local $2)
    )
    (i32.load offset=8
     (get_local $2)
    )
    (i32.const 32)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $2)
    (i32.const 16)
   )
  )
  (i32.eq
   (get_local $1)
   (i32.const 0)
  )
 )
 (func $_ZneRK11checksum160S1_ (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (set_local $2
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 16)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $2)
  )
  (i32.store offset=12
   (get_local $2)
   (get_local $0)
  )
  (i32.store offset=8
   (get_local $2)
   (get_local $1)
  )
  (set_local $1
   (call $memcmp
    (i32.load offset=12
     (get_local $2)
    )
    (i32.load offset=8
     (get_local $2)
    )
    (i32.const 32)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $2)
    (i32.const 16)
   )
  )
  (i32.ne
   (get_local $1)
   (i32.const 0)
  )
 )
 (func $transfer (param $0 i32) (param $1 i32) (param $2 i32) (result i32)
  (local $3 i32)
  (set_local $3
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 16)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $3)
  )
  (i32.store offset=12
   (get_local $3)
   (get_local $0)
  )
  (i32.store offset=8
   (get_local $3)
   (get_local $1)
  )
  (i32.store offset=4
   (get_local $3)
   (get_local $2)
  )
  (set_local $2
   (call $_ZN6platon5Token8transferEPcS1_i
    (call $_ZN6platon5TokenC2Ev
     (get_local $3)
    )
    (i32.load offset=12
     (get_local $3)
    )
    (i32.load offset=8
     (get_local $3)
    )
    (i32.load offset=4
     (get_local $3)
    )
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $3)
    (i32.const 16)
   )
  )
  (get_local $2)
 )
 (func $_ZN6platon5TokenC2Ev (param $0 i32) (result i32)
  (local $1 i32)
  (set_local $1
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 16)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $1)
  )
  (i32.store offset=12
   (get_local $1)
   (get_local $0)
  )
  (set_local $0
   (i32.load offset=12
    (get_local $1)
   )
  )
  (drop
   (call $_ZN6platon8ContractC2Ev
    (get_local $0)
   )
  )
  (i32.store
   (get_local $0)
   (i32.add
    (i32.const 52)
    (i32.const 8)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $1)
    (i32.const 16)
   )
  )
  (get_local $0)
 )
 (func $_ZN6platon5Token8transferEPcS1_i (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (result i32)
  (local $4 i32)
  (set_local $4
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 10032)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $4)
  )
  (i32.store offset=10028
   (get_local $4)
   (get_local $0)
  )
  (i32.store offset=10024
   (get_local $4)
   (get_local $1)
  )
  (i32.store offset=10020
   (get_local $4)
   (get_local $2)
  )
  (i32.store offset=10016
   (get_local $4)
   (get_local $3)
  )
  (set_local $3
   (i32.load offset=10028
    (get_local $4)
   )
  )
  (call $_ZN6platon7print_fIPcJS1_iEEEvPKcT_DpT0_
   (i32.const 16)
   (i32.load offset=10024
    (get_local $4)
   )
   (i32.load offset=10020
    (get_local $4)
   )
   (i32.load offset=10016
    (get_local $4)
   )
  )
  (i32.store offset=12
   (get_local $4)
   (i32.const 0)
  )
  (block $label$0
   (loop $label$1
    (br_if $label$0
     (i32.ge_u
      (i32.load offset=12
       (get_local $4)
      )
      (i32.const 10000)
     )
    )
    (i32.store8
     (i32.add
      (i32.add
       (get_local $4)
       (i32.const 16)
      )
      (i32.load offset=12
       (get_local $4)
      )
     )
     (i32.load offset=12
      (get_local $4)
     )
    )
    (i32.store offset=12
     (get_local $4)
     (i32.add
      (i32.load offset=12
       (get_local $4)
      )
      (i32.const 1)
     )
    )
    (br $label$1)
   )
  )
  (call $_ZN6platon5Token4testEv
   (get_local $3)
  )
  (drop
   (call $strlen
    (i32.load offset=10024
     (get_local $4)
    )
   )
  )
  (drop
   (call $atoi
    (i32.const 48)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $4)
    (i32.const 10032)
   )
  )
  (i32.const 88)
 )
 (func $_ZN6platon7print_fIPcJS1_iEEEvPKcT_DpT0_ (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32)
  (local $4 i32)
  (set_local $4
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 16)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $4)
  )
  (i32.store offset=12
   (get_local $4)
   (get_local $0)
  )
  (i32.store offset=8
   (get_local $4)
   (get_local $1)
  )
  (i32.store offset=4
   (get_local $4)
   (get_local $2)
  )
  (i32.store
   (get_local $4)
   (get_local $3)
  )
  (block $label$0
   (block $label$1
    (loop $label$2
     (br_if $label$1
      (i32.eqz
       (i32.shr_s
        (i32.shl
         (i32.load8_u
          (i32.load offset=12
           (get_local $4)
          )
         )
         (i32.const 24)
        )
        (i32.const 24)
       )
      )
     )
     (block $label$3
      (br_if $label$3
       (i32.ne
        (i32.shr_s
         (i32.shl
          (i32.load8_u
           (i32.load offset=12
            (get_local $4)
           )
          )
          (i32.const 24)
         )
         (i32.const 24)
        )
        (i32.const 37)
       )
      )
      (call $_ZN6platon5printEPKc
       (i32.load offset=8
        (get_local $4)
       )
      )
      (call $_ZN6platon7print_fIPcJiEEEvPKcT_DpT0_
       (i32.add
        (i32.load offset=12
         (get_local $4)
        )
        (i32.const 1)
       )
       (i32.load offset=4
        (get_local $4)
       )
       (i32.load
        (get_local $4)
       )
      )
      (br $label$0)
     )
     (call $prints_l
      (i32.load offset=12
       (get_local $4)
      )
      (i32.const 1)
     )
     (i32.store offset=12
      (get_local $4)
      (i32.add
       (i32.load offset=12
        (get_local $4)
       )
       (i32.const 1)
      )
     )
     (br $label$2)
    )
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $4)
    (i32.const 16)
   )
  )
 )
 (func $_ZN6platon5Token4testEv (param $0 i32)
  (local $1 i32)
  (set_local $1
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 10032)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $1)
  )
  (i32.store offset=10028
   (get_local $1)
   (get_local $0)
  )
  (i32.store offset=12
   (get_local $1)
   (i32.const 0)
  )
  (block $label$0
   (loop $label$1
    (br_if $label$0
     (i32.ge_u
      (i32.load offset=12
       (get_local $1)
      )
      (i32.const 10000)
     )
    )
    (i32.store8
     (i32.add
      (i32.add
       (get_local $1)
       (i32.const 16)
      )
      (i32.load offset=12
       (get_local $1)
      )
     )
     (i32.load offset=12
      (get_local $1)
     )
    )
    (i32.store offset=12
     (get_local $1)
     (i32.add
      (i32.load offset=12
       (get_local $1)
      )
      (i32.const 1)
     )
    )
    (br $label$1)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $1)
    (i32.const 10032)
   )
  )
 )
 (func $_ZN6platon5printEPKc (param $0 i32)
  (local $1 i32)
  (set_local $1
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 16)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $1)
  )
  (i32.store offset=12
   (get_local $1)
   (get_local $0)
  )
  (call $prints
   (i32.load offset=12
    (get_local $1)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $1)
    (i32.const 16)
   )
  )
 )
 (func $_ZN6platon7print_fIPcJiEEEvPKcT_DpT0_ (param $0 i32) (param $1 i32) (param $2 i32)
  (local $3 i32)
  (set_local $3
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 16)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $3)
  )
  (i32.store offset=12
   (get_local $3)
   (get_local $0)
  )
  (i32.store offset=8
   (get_local $3)
   (get_local $1)
  )
  (i32.store offset=4
   (get_local $3)
   (get_local $2)
  )
  (block $label$0
   (block $label$1
    (loop $label$2
     (br_if $label$1
      (i32.eqz
       (i32.shr_s
        (i32.shl
         (i32.load8_u
          (i32.load offset=12
           (get_local $3)
          )
         )
         (i32.const 24)
        )
        (i32.const 24)
       )
      )
     )
     (block $label$3
      (br_if $label$3
       (i32.ne
        (i32.shr_s
         (i32.shl
          (i32.load8_u
           (i32.load offset=12
            (get_local $3)
           )
          )
          (i32.const 24)
         )
         (i32.const 24)
        )
        (i32.const 37)
       )
      )
      (call $_ZN6platon5printEPKc
       (i32.load offset=8
        (get_local $3)
       )
      )
      (call $_ZN6platon7print_fIiJEEEvPKcT_DpT0_
       (i32.add
        (i32.load offset=12
         (get_local $3)
        )
        (i32.const 1)
       )
       (i32.load offset=4
        (get_local $3)
       )
      )
      (br $label$0)
     )
     (call $prints_l
      (i32.load offset=12
       (get_local $3)
      )
      (i32.const 1)
     )
     (i32.store offset=12
      (get_local $3)
      (i32.add
       (i32.load offset=12
        (get_local $3)
       )
       (i32.const 1)
      )
     )
     (br $label$2)
    )
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $3)
    (i32.const 16)
   )
  )
 )
 (func $_ZN6platon7print_fIiJEEEvPKcT_DpT0_ (param $0 i32) (param $1 i32)
  (local $2 i32)
  (set_local $2
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 16)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $2)
  )
  (i32.store offset=12
   (get_local $2)
   (get_local $0)
  )
  (i32.store offset=8
   (get_local $2)
   (get_local $1)
  )
  (block $label$0
   (block $label$1
    (loop $label$2
     (br_if $label$1
      (i32.eqz
       (i32.shr_s
        (i32.shl
         (i32.load8_u
          (i32.load offset=12
           (get_local $2)
          )
         )
         (i32.const 24)
        )
        (i32.const 24)
       )
      )
     )
     (block $label$3
      (br_if $label$3
       (i32.ne
        (i32.shr_s
         (i32.shl
          (i32.load8_u
           (i32.load offset=12
            (get_local $2)
           )
          )
          (i32.const 24)
         )
         (i32.const 24)
        )
        (i32.const 37)
       )
      )
      (call $_ZN6platon5printEi
       (i32.load offset=8
        (get_local $2)
       )
      )
      (call $_ZN6platon7print_fEPKc
       (i32.add
        (i32.load offset=12
         (get_local $2)
        )
        (i32.const 1)
       )
      )
      (br $label$0)
     )
     (call $prints_l
      (i32.load offset=12
       (get_local $2)
      )
      (i32.const 1)
     )
     (i32.store offset=12
      (get_local $2)
      (i32.add
       (i32.load offset=12
        (get_local $2)
       )
       (i32.const 1)
      )
     )
     (br $label$2)
    )
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $2)
    (i32.const 16)
   )
  )
 )
 (func $_ZN6platon5printEi (param $0 i32)
  (local $1 i32)
  (set_local $1
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 16)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $1)
  )
  (i32.store offset=12
   (get_local $1)
   (get_local $0)
  )
  (call $printi
   (i64.extend_s/i32
    (i32.load offset=12
     (get_local $1)
    )
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $1)
    (i32.const 16)
   )
  )
 )
 (func $_ZN6platon7print_fEPKc (param $0 i32)
  (local $1 i32)
  (set_local $1
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 16)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $1)
  )
  (i32.store offset=12
   (get_local $1)
   (get_local $0)
  )
  (call $prints
   (i32.load offset=12
    (get_local $1)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $1)
    (i32.const 16)
   )
  )
 )
 (func $_ZN6platon8ContractC2Ev (param $0 i32) (result i32)
  (local $1 i32)
  (set_local $1
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 16)
   )
  )
  (i32.store offset=12
   (get_local $1)
   (get_local $0)
  )
  (set_local $1
   (i32.load offset=12
    (get_local $1)
   )
  )
  (i32.store
   (get_local $1)
   (i32.add
    (i32.const 64)
    (i32.const 8)
   )
  )
  (get_local $1)
 )
 (func $_ZN6platon5Token4initEv (type $FUNCSIG$vi) (param $0 i32)
  (i32.store offset=12
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 16)
   )
   (get_local $0)
  )
 )
 (func $atoi (param $0 i32) (result i32)
  (local $1 i32)
  (local $2 i32)
  (set_local $2
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 16)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $2)
  )
  (i32.store offset=12
   (get_local $2)
   (get_local $0)
  )
  (i32.store offset=8
   (get_local $2)
   (i32.const 0)
  )
  (i32.store offset=4
   (get_local $2)
   (i32.const 0)
  )
  (block $label$0
   (loop $label$1
    (br_if $label$0
     (i32.eqz
      (call $__isspace.502
       (i32.shr_s
        (i32.shl
         (i32.load8_u
          (i32.load offset=12
           (get_local $2)
          )
         )
         (i32.const 24)
        )
        (i32.const 24)
       )
      )
     )
    )
    (i32.store offset=12
     (get_local $2)
     (i32.add
      (i32.load offset=12
       (get_local $2)
      )
      (i32.const 1)
     )
    )
    (br $label$1)
   )
  )
  (set_local $0
   (i32.load8_s
    (i32.load offset=12
     (get_local $2)
    )
   )
  )
  (block $label$2
   (block $label$3
    (br_if $label$3
     (i32.eq
      (get_local $0)
      (i32.const 43)
     )
    )
    (br_if $label$2
     (i32.ne
      (get_local $0)
      (i32.const 45)
     )
    )
    (i32.store offset=4
     (get_local $2)
     (i32.const 1)
    )
   )
   (i32.store offset=12
    (get_local $2)
    (i32.add
     (i32.load offset=12
      (get_local $2)
     )
     (i32.const 1)
    )
   )
  )
  (block $label$4
   (loop $label$5
    (br_if $label$4
     (i32.ge_u
      (i32.sub
       (i32.shr_s
        (i32.shl
         (i32.load8_u
          (i32.load offset=12
           (get_local $2)
          )
         )
         (i32.const 24)
        )
        (i32.const 24)
       )
       (i32.const 48)
      )
      (i32.const 10)
     )
    )
    (set_local $1
     (i32.load offset=8
      (get_local $2)
     )
    )
    (set_local $0
     (i32.load offset=12
      (get_local $2)
     )
    )
    (i32.store offset=12
     (get_local $2)
     (i32.add
      (get_local $0)
      (i32.const 1)
     )
    )
    (i32.store offset=8
     (get_local $2)
     (i32.sub
      (i32.mul
       (get_local $1)
       (i32.const 10)
      )
      (i32.sub
       (i32.shr_s
        (i32.shl
         (i32.load8_u
          (get_local $0)
         )
         (i32.const 24)
        )
        (i32.const 24)
       )
       (i32.const 48)
      )
     )
    )
    (br $label$5)
   )
  )
  (block $label$6
   (block $label$7
    (br_if $label$7
     (i32.eqz
      (i32.load offset=4
       (get_local $2)
      )
     )
    )
    (set_local $0
     (i32.load offset=8
      (get_local $2)
     )
    )
    (br $label$6)
   )
   (set_local $0
    (i32.sub
     (i32.const 0)
     (i32.load offset=8
      (get_local $2)
     )
    )
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $2)
    (i32.const 16)
   )
  )
  (get_local $0)
 )
 (func $__isspace.502 (param $0 i32) (result i32)
  (local $1 i32)
  (local $2 i32)
  (set_local $2
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 16)
   )
  )
  (set_local $1
   (i32.const 1)
  )
  (i32.store offset=12
   (get_local $2)
   (get_local $0)
  )
  (block $label$0
   (br_if $label$0
    (i32.eq
     (i32.load offset=12
      (get_local $2)
     )
     (i32.const 32)
    )
   )
   (set_local $1
    (i32.lt_u
     (i32.sub
      (i32.load offset=12
       (get_local $2)
      )
      (i32.const 9)
     )
     (i32.const 5)
    )
   )
  )
  (i32.and
   (get_local $1)
   (i32.const 1)
  )
 )
 (func $memcmp (param $0 i32) (param $1 i32) (param $2 i32) (result i32)
  (local $3 i32)
  (set_local $3
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 32)
   )
  )
  (i32.store offset=28
   (get_local $3)
   (get_local $0)
  )
  (i32.store offset=24
   (get_local $3)
   (get_local $1)
  )
  (i32.store offset=20
   (get_local $3)
   (get_local $2)
  )
  (i32.store offset=16
   (get_local $3)
   (i32.load offset=28
    (get_local $3)
   )
  )
  (i32.store offset=12
   (get_local $3)
   (i32.load offset=24
    (get_local $3)
   )
  )
  (loop $label$0
   (set_local $2
    (i32.const 0)
   )
   (block $label$1
    (br_if $label$1
     (i32.eqz
      (i32.load offset=20
       (get_local $3)
      )
     )
    )
    (set_local $2
     (i32.eq
      (i32.and
       (i32.load8_u
        (i32.load offset=16
         (get_local $3)
        )
       )
       (i32.const 255)
      )
      (i32.and
       (i32.load8_u
        (i32.load offset=12
         (get_local $3)
        )
       )
       (i32.const 255)
      )
     )
    )
   )
   (block $label$2
    (br_if $label$2
     (i32.eqz
      (i32.and
       (get_local $2)
       (i32.const 1)
      )
     )
    )
    (i32.store offset=20
     (get_local $3)
     (i32.add
      (i32.load offset=20
       (get_local $3)
      )
      (i32.const -1)
     )
    )
    (i32.store offset=16
     (get_local $3)
     (i32.add
      (i32.load offset=16
       (get_local $3)
      )
      (i32.const 1)
     )
    )
    (i32.store offset=12
     (get_local $3)
     (i32.add
      (i32.load offset=12
       (get_local $3)
      )
      (i32.const 1)
     )
    )
    (br $label$0)
   )
  )
  (block $label$3
   (block $label$4
    (br_if $label$4
     (i32.eqz
      (i32.load offset=20
       (get_local $3)
      )
     )
    )
    (set_local $3
     (i32.sub
      (i32.and
       (i32.load8_u
        (i32.load offset=16
         (get_local $3)
        )
       )
       (i32.const 255)
      )
      (i32.and
       (i32.load8_u
        (i32.load offset=12
         (get_local $3)
        )
       )
       (i32.const 255)
      )
     )
    )
    (br $label$3)
   )
   (set_local $3
    (i32.const 0)
   )
  )
  (get_local $3)
 )
 (func $strlen (param $0 i32) (result i32)
  (local $1 i32)
  (set_local $1
   (i32.sub
    (i32.load offset=4
     (i32.const 0)
    )
    (i32.const 16)
   )
  )
  (i32.store offset=8
   (get_local $1)
   (get_local $0)
  )
  (i32.store offset=4
   (get_local $1)
   (i32.load offset=8
    (get_local $1)
   )
  )
  (block $label$0
   (block $label$1
    (loop $label$2
     (br_if $label$1
      (i32.eqz
       (i32.and
        (i32.load offset=8
         (get_local $1)
        )
        (i32.const 3)
       )
      )
     )
     (block $label$3
      (br_if $label$3
       (i32.ne
        (i32.and
         (i32.load8_u
          (i32.load offset=8
           (get_local $1)
          )
         )
         (i32.const 255)
        )
        (i32.and
         (i32.const 0)
         (i32.const 255)
        )
       )
      )
      (i32.store offset=12
       (get_local $1)
       (i32.sub
        (i32.load offset=8
         (get_local $1)
        )
        (i32.load offset=4
         (get_local $1)
        )
       )
      )
      (br $label$0)
     )
     (i32.store offset=8
      (get_local $1)
      (i32.add
       (i32.load offset=8
        (get_local $1)
       )
       (i32.const 1)
      )
     )
     (br $label$2)
    )
   )
   (i32.store
    (get_local $1)
    (i32.load offset=8
     (get_local $1)
    )
   )
   (block $label$4
    (loop $label$5
     (br_if $label$4
      (i32.ne
       (i32.and
        (i32.and
         (i32.sub
          (i32.load
           (i32.load
            (get_local $1)
           )
          )
          (i32.const 16843009)
         )
         (i32.xor
          (i32.load
           (i32.load
            (get_local $1)
           )
          )
          (i32.const -1)
         )
        )
        (i32.const -2139062144)
       )
       (i32.const 0)
      )
     )
     (i32.store
      (get_local $1)
      (i32.add
       (i32.load
        (get_local $1)
       )
       (i32.const 4)
      )
     )
     (br $label$5)
    )
   )
   (i32.store offset=8
    (get_local $1)
    (i32.load
     (get_local $1)
    )
   )
   (block $label$6
    (loop $label$7
     (br_if $label$6
      (i32.eq
       (i32.and
        (i32.load8_u
         (i32.load offset=8
          (get_local $1)
         )
        )
        (i32.const 255)
       )
       (i32.and
        (i32.const 0)
        (i32.const 255)
       )
      )
     )
     (i32.store offset=8
      (get_local $1)
      (i32.add
       (i32.load offset=8
        (get_local $1)
       )
       (i32.const 1)
      )
     )
     (br $label$7)
    )
   )
   (i32.store offset=12
    (get_local $1)
    (i32.sub
     (i32.load offset=8
      (get_local $1)
     )
     (i32.load offset=4
      (get_local $1)
     )
    )
   )
  )
  (i32.load offset=12
   (get_local $1)
  )
 )
 (func $__wasm_nullptr (type $FUNCSIG$v)
  (unreachable)
 )
 (func $__importThunk___cxa_pure_virtual (type $FUNCSIG$v)
  (call $__cxa_pure_virtual)
 )
)
