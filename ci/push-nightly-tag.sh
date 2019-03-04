#!/bin/sh

echo "Starting nightly build deploy script..."
MY_TAG="$( git describe --exact-match "$(git rev-parse HEAD)" 2>/dev/null )"
if [ -z "$MY_TAG" ] ; then
  echo "Tag for last commit is not found, going to try to push new nighlty build tag..."

  git config --global user.email "travis@travis-ci.org"
  git config --global user.name "Travis CI"

  NEW_TAG="Nightly-$TRAVIS_BRANCH-$(date +%Y-%m-%d)-$(git rev-parse --short HEAD)"
  git tag -a $NEW_TAG -m "Nightly Build Tag $NEW_TAG"

  echo "New generated nightly build tag: $NEW_TAG"

  git remote add origin-repo https://${GITHUB_TOKEN}@github.com/liamg/aminal.git > /dev/null 2>&1
  git push origin-repo $NEW_TAG
else
  echo "Skipping nighly build tag generation. Last commit tag found:$MY_TAG"
fi
