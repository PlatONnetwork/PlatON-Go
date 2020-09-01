class Gas {
  const char* name_ = nullptr;
  uint64_t gas_;

 public:
  Gas(const char *name) : name_(name), gas_(platon_gas()) {}
  ~Gas()
  {
    emit();
  }
  PLATON_EVENT1(GasUsed, const std::string &, uint64_t)
  void Reset(const char *name) {
    emit();
    name_ = name;
    gas_ = platon_gas();
  }
  void emit() {
    uint64_t cost = gas_ - platon_gas();
    PLATON_EMIT_EVENT1(GasUsed, name_, cost);
  }
};
