package evolve

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/base"
)

/**
 * <p>An evolutionary operator is a function that takes a population of
 * candidates as an argument and returns a new population that is the
 * result of applying a transformation to the original population.</p>
 * <p><strong>An implementation of this class must not modify any of
 * the selected candidate objects passed in.</strong>  Doing so will
 * affect the correct operation of the {@link EvolutionEngine}.  Instead
 * the operator should create and return new candidate objects.  The
 * operator is not required to create copies of unmodified individuals
 * (for efficiency these may be returned directly).</p>
 * @param <T> The type of evolvable entity that this operator accepts.
 * @author Daniel Dyer
 */
type EvolutionaryOperator interface {
	/**
	 * <p>Apply the operation to each entry in the list of selected
	 * candidates.  It is important to note that this method operates on
	 * the list of candidates returned by the selection strategy and not
	 * on the current population.  Each entry in the list (not each
	 * individual - the list may contain the same individual more than
	 * once) must be operated on exactly once.</p>
	 *
	 * <p>Implementing classes should not assume any particular ordering
	 * (or lack of ordering) for the selection.  If ordering or
	 * shuffling is required, it should be performed by the implementing
	 * class.  The implementation should not re-order the list provided
	 * but instead should make a copy of the list and re-order that.
	 * The ordering of the selection should be totally irrelevant for
	 * operators that process each candidate in isolation, such as mutation.
	 * It should only be an issue for operators, such as cross-over, that
	 * deal with multiple candidates in a single operation.</p>
	 * <p><strong>The operator must not modify any of the candidates passed
	 * in</strong>.  Instead it should return a list that contains evolved
	 * copies of those candidates (umodified candidates can be included in
	 * the results without having to be copied).</p>
	 * @param selectedCandidates The individuals to evolve.
	 * @param rng A source of randomness for stochastic operators (most
	 * operators will be stochastic).
	 * @return The evolved individuals.
	 */
	Apply(selectedCandidates []base.Candidate, rng *rand.Rand) []base.Candidate
}
