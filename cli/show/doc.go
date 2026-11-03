/*
Package show implements the CLI args filtering which kinds of capture targets
are offered in the “Containers” selector of the extcap configuration dialog.
Wireshark passes these args from the “Standard” tab of the dialog when
reloading the selector:

  - “--showpods” shows Kubernetes pods.
  - “--showcontainers” shows stand-alone containers, that is, containers not
    belonging to a pod.
  - “--showprocs” shows processes with their own virtual IP stacks that
    aren't containers.
  - “--showemptynetns” shows IP stacks without any attached processes, such
    as bind-mounted network namespaces.

All args default to hiding the respective kind of capture target; use
[Target] to check whether a particular capture target is to be shown.
*/
package show
