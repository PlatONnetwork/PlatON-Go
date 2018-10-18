(module
 (type $FUNCSIG$vi (func (param i32)))
 (type $FUNCSIG$vj (func (param i64)))
 (type $FUNCSIG$vii (func (param i32 i32)))
 (type $FUNCSIG$v (func))
 (import "env" "printi" (func $printi (param i64)))
 (import "env" "prints" (func $prints (param i32)))
 (import "env" "prints_l" (func $prints_l (param i32 i32)))
 (table 2 2 anyfunc)
 (elem (i32.const 0) $__wasm_nullptr $_ZN6platon5Token4initEv)
 (memory $0 2)
 (data (i32.const 4) "@\90\01\00")
 (data (i32.const 12) "\00\00\00\00\00\00\00\00\01\00\00\00")
 (data (i32.const 32) "from:% to:% asset: % \n\00")
 (export "memory" (memory $0))
 (export "_ZeqRK11checksum256S1_" (func $_ZeqRK11checksum256S1_))
 (export "_ZeqRK11checksum160S1_" (func $_ZeqRK11checksum160S1_))
 (export "_ZneRK11checksum160S1_" (func $_ZneRK11checksum160S1_))
 (export "transfer" (func $transfer))
 (export "init" (func $init))
 (export "memcmp" (func $memcmp))
 (func $_ZeqRK11checksum256S1_ (param $0 i32) (param $1 i32) (result i32)
  (i32.eqz
   (call $memcmp
    (get_local $0)
    (get_local $1)
    (i32.const 32)
   )
  )
 )
 (func $_ZeqRK11checksum160S1_ (param $0 i32) (param $1 i32) (result i32)
  (i32.eqz
   (call $memcmp
    (get_local $0)
    (get_local $1)
    (i32.const 32)
   )
  )
 )
 (func $_ZneRK11checksum160S1_ (param $0 i32) (param $1 i32) (result i32)
  (i32.ne
   (call $memcmp
    (get_local $0)
    (get_local $1)
    (i32.const 32)
   )
   (i32.const 0)
  )
 )
 (func $transfer (param $0 i32) (param $1 i32) (param $2 i32)
  (local $3 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $3
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 16)
    )
   )
  )
  (i32.store offset=8
   (get_local $3)
   (i32.const 20)
  )
  (call $_ZN6platon5Token8transferEPcS1_i
   (i32.add
    (get_local $3)
    (i32.const 8)
   )
   (get_local $0)
   (get_local $1)
   (get_local $2)
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $3)
    (i32.const 16)
   )
  )
 )
 (func $_ZN6platon5Token8transferEPcS1_i (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32)
  (local $4 i32)
  (local $5 i32)
  (set_local $5
   (i32.const 32)
  )
  (block $label$0
   (br_if $label$0
    (i32.eqz
     (tee_local $4
      (i32.load8_u
       (i32.const 32)
      )
     )
    )
   )
   (block $label$1
    (loop $label$2
     (br_if $label$1
      (i32.eq
       (get_local $4)
       (i32.const 37)
      )
     )
     (call $prints_l
      (get_local $5)
      (i32.const 1)
     )
     (br_if $label$0
      (i32.eqz
       (tee_local $4
        (i32.load8_u
         (tee_local $5
          (i32.add
           (get_local $5)
           (i32.const 1)
          )
         )
        )
       )
      )
     )
     (br $label$2)
    )
   )
   (call $prints
    (get_local $1)
   )
   (br_if $label$0
    (i32.eqz
     (tee_local $4
      (i32.load8_u
       (tee_local $5
        (i32.add
         (get_local $5)
         (i32.const 1)
        )
       )
      )
     )
    )
   )
   (block $label$3
    (loop $label$4
     (br_if $label$3
      (i32.eq
       (get_local $4)
       (i32.const 37)
      )
     )
     (call $prints_l
      (get_local $5)
      (i32.const 1)
     )
     (br_if $label$0
      (i32.eqz
       (tee_local $4
        (i32.load8_u
         (tee_local $5
          (i32.add
           (get_local $5)
           (i32.const 1)
          )
         )
        )
       )
      )
     )
     (br $label$4)
    )
   )
   (call $prints
    (get_local $2)
   )
   (br_if $label$0
    (i32.eqz
     (tee_local $4
      (i32.load8_u
       (tee_local $5
        (i32.add
         (get_local $5)
         (i32.const 1)
        )
       )
      )
     )
    )
   )
   (block $label$5
    (loop $label$6
     (br_if $label$5
      (i32.eq
       (get_local $4)
       (i32.const 37)
      )
     )
     (call $prints_l
      (get_local $5)
      (i32.const 1)
     )
     (br_if $label$0
      (i32.eqz
       (tee_local $4
        (i32.load8_u
         (tee_local $5
          (i32.add
           (get_local $5)
           (i32.const 1)
          )
         )
        )
       )
      )
     )
     (br $label$6)
    )
   )
   (call $printi
    (i64.extend_s/i32
     (get_local $3)
    )
   )
   (call $prints
    (i32.add
     (get_local $5)
     (i32.const 1)
    )
   )
  )
 )
 (func $_ZN6platon5Token4initEv (type $FUNCSIG$vi) (param $0 i32)
 )
 (func $init
 )
 (func $memcmp (param $0 i32) (param $1 i32) (param $2 i32) (result i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (set_local $5
   (i32.const 0)
  )
  (block $label$0
   (br_if $label$0
    (i32.eqz
     (get_local $2)
    )
   )
   (block $label$1
    (loop $label$2
     (br_if $label$1
      (i32.ne
       (tee_local $3
        (i32.load8_u
         (get_local $0)
        )
       )
       (tee_local $4
        (i32.load8_u
         (get_local $1)
        )
       )
      )
     )
     (set_local $1
      (i32.add
       (get_local $1)
       (i32.const 1)
      )
     )
     (set_local $0
      (i32.add
       (get_local $0)
       (i32.const 1)
      )
     )
     (br_if $label$2
      (tee_local $2
       (i32.add
        (get_local $2)
        (i32.const -1)
       )
      )
     )
     (br $label$0)
    )
   )
   (set_local $5
    (i32.sub
     (get_local $3)
     (get_local $4)
    )
   )
  )
  (get_local $5)
 )
 (func $__wasm_nullptr (type $FUNCSIG$v)
  (unreachable)
 )
)
