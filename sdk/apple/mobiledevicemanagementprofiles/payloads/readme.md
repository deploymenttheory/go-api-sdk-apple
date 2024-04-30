
# Use of pointers for boolean fields in Go structs

Go structs for declatative management profiles uses pointers for boolean fields (*bool) rather than (bool) as it allows for field three states:

true: Explicitly enable the setting.
false: Explicitly disable the setting.
nil: Use the default behavior, which is typically not setting the policy, thereby leaving the configuration at its system or previously configured value.