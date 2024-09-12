#!/bin/bash

# Generate buffalo resources
buffalo generate resource item --skip-templates
buffalo generate resource customer --skip-templates --skip-migration
buffalo generate resource business-details --skip-templates --skip-migration
buffalo generate resource payment-details --skip-templates --skip-migration
buffalo generate resource invoice --skip-templates --skip-migration
buffalo g auth