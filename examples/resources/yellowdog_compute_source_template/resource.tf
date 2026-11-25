# A compute source template describes one place instances can be provisioned:
# a cloud provider, account, region, network and instance type. Compute
# requirement templates draw on them. Set exactly one of the *_compute_source
# attributes of source.
#
# credential names a keyring credential as "<keyring name>/<credential name>".

# AWS instances, launched individually.
resource "yellowdog_compute_source_template" "aws_instances" {
  namespace   = "example-namespace"
  description = "On-demand t3a instances in eu-west-2"

  source = {
    aws_instances_compute_source = {
      name              = "aws-eu-west-2"
      credential        = "example-keyring/aws"
      region            = "eu-west-2"
      image_id          = "ami-09dbc7ce74870d573"
      instance_type     = "t3a.micro"
      security_group_id = "sg-0123456789abcdef0"
      subnet_id         = "subnet-0123456789abcdef0"
      assign_public_ip  = true
      limit             = 10

      instance_tags = {
        team = "research"
      }
    }
  }

  # User attributes, for dynamic compute requirement templates to constrain
  # and rank sources by. Each must already be defined in the account.
  attributes = [
    {
      string_attribute_value = {
        attribute = "user.tier"
        value     = "production"
      }
    },
    {
      numeric_attribute_value = {
        attribute = "user.priority"
        value     = 10
      }
    },
  ]
}

# An AWS EC2 Fleet, which can mix instance types and spot and on-demand
# capacity.
resource "yellowdog_compute_source_template" "aws_fleet" {
  namespace = "example-namespace"

  source = {
    aws_fleet_compute_source = {
      name              = "aws-fleet-eu-west-2"
      credential        = "example-keyring/aws"
      region            = "eu-west-2"
      image_id          = "ami-09dbc7ce74870d573"
      instance_type     = "t3a.micro"
      security_group_id = "sg-0123456789abcdef0"
      subnet_id         = "subnet-0123456789abcdef0"
      purchase_option   = "SPOT_THEN_ON_DEMAND"
      limit             = 50

      instance_overrides = [
        { instance_type = "t3a.small" },
        { instance_type = "t3.small" },
      ]

      spot_options = {
        allocation_strategy = "PRICE_CAPACITY_OPTIMIZED"
      }
    }
  }
}

# Azure virtual machines, launched individually.
resource "yellowdog_compute_source_template" "azure_instances" {
  namespace = "example-namespace"

  source = {
    azure_instances_compute_source = {
      name                        = "azure-uksouth"
      credential                  = "example-keyring/azure"
      admin_user_credential       = "example-keyring/azure-admin"
      region                      = "uksouth"
      image_id                    = "market:Canonical:0001-com-ubuntu-server-jammy:22_04-lts:latest"
      vm_size                     = "Standard_B1s"
      network_resource_group_name = "example-network-rg"
      network_name                = "example-vnet"
      subnet_name                 = "default"
      limit                       = 10
    }
  }
}

# An Azure virtual machine scale set.
resource "yellowdog_compute_source_template" "azure_scale_set" {
  namespace = "example-namespace"

  source = {
    azure_scale_set_compute_source = {
      name                        = "azure-scale-set-uksouth"
      credential                  = "example-keyring/azure"
      admin_user_credential       = "example-keyring/azure-admin"
      region                      = "uksouth"
      image_id                    = "market:Canonical:0001-com-ubuntu-server-jammy:22_04-lts:latest"
      vm_size                     = "Standard_B1s"
      network_resource_group_name = "example-network-rg"
      network_name                = "example-vnet"
      subnet_name                 = "default"
      use_spot                    = true
      limit                       = 50
    }
  }
}

# Google Compute Engine instances, launched individually.
resource "yellowdog_compute_source_template" "gce_instances" {
  namespace = "example-namespace"

  source = {
    gce_instances_compute_source = {
      name         = "gce-europe-west2"
      credential   = "example-keyring/gcp"
      project      = "example-project"
      region       = "europe-west2"
      zone         = "europe-west2-a"
      image        = "projects/ubuntu-os-cloud/global/images/ubuntu-2204-jammy-v20240904"
      machine_type = "e2-micro"
      network      = "default"
      subnetwork   = "default"
      limit        = 10
    }
  }
}

# A Google Compute Engine managed instance group.
resource "yellowdog_compute_source_template" "gce_instance_group" {
  namespace = "example-namespace"

  source = {
    gce_instance_group_compute_source = {
      name                      = "gce-group-europe-west2"
      credential                = "example-keyring/gcp"
      project                   = "example-project"
      region                    = "europe-west2"
      image                     = "projects/ubuntu-os-cloud/global/images/ubuntu-2204-jammy-v20240904"
      machine_type              = "e2-micro"
      target_distribution_shape = "BALANCED"
      spot                      = true
      limit                     = 50
    }
  }
}

# Oracle Cloud Infrastructure instances, launched individually.
resource "yellowdog_compute_source_template" "oci_instances" {
  namespace = "example-namespace"

  source = {
    oci_instances_compute_source = {
      name           = "oci-uk-london-1"
      credential     = "example-keyring/oci"
      region         = "uk-london-1"
      compartment_id = "ocid1.compartment.oc1..aaaaaaaaexampleexampleexampleexampleexampleexampleexample"
      subnet_id      = "ocid1.subnet.oc1.uk-london-1.aaaaaaaaexampleexampleexampleexampleexampleexampleexample"
      image_id       = "ocid1.image.oc1.uk-london-1.aaaaaaaaxwu543plnc7ulnexmpxz6eifxk7eqxnnbwxtnzsxpckw7ogfesrq"
      shape          = "VM.Standard.E5.Flex"
      # Flex shapes need their OCPU count and memory (GB).
      flex_ocpus = 1
      flex_ram   = 8
      limit      = 10
    }
  }
}

# An Oracle Cloud Infrastructure instance pool.
resource "yellowdog_compute_source_template" "oci_instance_pool" {
  namespace = "example-namespace"

  source = {
    oci_instance_pool_compute_source = {
      name           = "oci-pool-uk-london-1"
      credential     = "example-keyring/oci"
      region         = "uk-london-1"
      compartment_id = "ocid1.compartment.oc1..aaaaaaaaexampleexampleexampleexampleexampleexampleexample"
      subnet_id      = "ocid1.subnet.oc1.uk-london-1.aaaaaaaaexampleexampleexampleexampleexampleexampleexample"
      image_id       = "ocid1.image.oc1.uk-london-1.aaaaaaaaxwu543plnc7ulnexmpxz6eifxk7eqxnnbwxtnzsxpckw7ogfesrq"
      shape          = "VM.Standard.E5.Flex"
      flex_ocpus     = 1
      flex_ram       = 8
      limit          = 50
    }
  }
}

# A simulated source, which provisions no real instances: useful for trying
# out templates and workflows without a cloud account or cost.
resource "yellowdog_compute_source_template" "simulator" {
  namespace = "example-namespace"

  source = {
    simulator_compute_source = {
      name                          = "simulator"
      region                        = "sim-region-1"
      instance_type                 = "sim-small"
      image_id                      = "sim-image-1"
      implicit_capacity             = 100
      instance_startup_time_seconds = 30
    }
  }
}
