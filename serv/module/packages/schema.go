package packages

type PackagesObjectStruct struct {
	Id            int
	Name          string
	Price         int
	Duration      int
	Duration_type string
	Expired_date  string
}

type PackageGroupDenyContentStruct struct {
	Id                       int
	Billing_package_group_id int
	Content_id               string
	Name                     string
	Pk_id                    int
}
