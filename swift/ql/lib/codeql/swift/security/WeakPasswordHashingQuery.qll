/**
 * Provides a taint tracking configuration to find use of inappropriate
 * cryptographic hashing algorithms on passwords.
 */

import swift
import codeql.swift.security.SensitiveExprs
import codeql.swift.dataflow.DataFlow
import codeql.swift.dataflow.TaintTracking
import codeql.swift.security.WeakPasswordHashingExtensions

/**
 * A taint tracking configuration from password expressions to inappropriate
 * hashing sinks.
 */
module WeakHashingPasswordConfig implements DataFlow::ConfigSig {
  predicate isSource(DataFlow::Node node) { node.asExpr() instanceof PasswordExpr }

  predicate isSink(DataFlow::Node node) { node instanceof WeakPasswordHashingSink }

  predicate isBarrier(DataFlow::Node node) { node instanceof WeakPasswordHashingBarrier }

  predicate isBarrierIn(DataFlow::Node node) {
    // make sources barriers so that we only report the closest instance
    isSource(node)
  }

  predicate isBarrierOut(DataFlow::Node node) {
    // make sinks barriers so that we only report the closest instance
    isSink(node)
  }

  predicate isAdditionalFlowStep(DataFlow::Node nodeFrom, DataFlow::Node nodeTo) {
    any(WeakPasswordHashingAdditionalFlowStep s).step(nodeFrom, nodeTo)
  }
}

module WeakHashingFlow = TaintTracking::Global<WeakHashingPasswordConfig>;
